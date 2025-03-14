import pandas as pd
import pickle
from sklearn.preprocessing import LabelEncoder
import os
from tqdm import tqdm
import torch
import torch.nn as nn
import torch.optim as optim
from sklearn.model_selection import train_test_split
import numpy as np

# Function to load formations data
def load_formations_data(directory="data/formations"):
    """Load all formation CSV files, excluding 'players' column."""
    all_formations = []
    for filename in tqdm(os.listdir(directory), desc="Loading formations CSVs"):
        if filename.startswith("formations_") and filename.endswith(".csv"):
            file_path = os.path.join(directory, filename)
            df = pd.read_csv(file_path)[['gameid', 'team', 'formation']]
            all_formations.append(df)
    return pd.concat(all_formations, ignore_index=True)

# Load data
formations_df = load_formations_data()
results_df = pd.read_parquet("data/results/games.parquet")

# Clean team names
def clean_team_name(team_name):
    if pd.isna(team_name):
        return None
    return str(team_name).replace(" (England)", "").strip()

results_df['home_clean'] = results_df['home'].apply(clean_team_name)
results_df['away_clean'] = results_df['away'].apply(clean_team_name)
formations_df['team_clean'] = formations_df['team'].apply(clean_team_name)

# Clean formation strings
def clean_formation(formation):
    return formation.strip().strip("'\"")

formations_df['formation'] = formations_df['formation'].apply(clean_formation)

# Merge for home and away teams
data_home = results_df.merge(formations_df, 
                             left_on=["home_clean"], 
                             right_on=["team_clean"], 
                             how="inner")
data_away = results_df.merge(formations_df, 
                             left_on=["away_clean"], 
                             right_on=["team_clean"], 
                             how="inner")

# Rename formation columns
data_home = data_home.rename(columns={'formation': 'home_formation'})
data_away = data_away.rename(columns={'formation': 'away_formation'})

# Combine into one dataframe
data = pd.merge(data_home[['date', 'home_clean', 'away_clean', 'gh', 'ga', 'home_formation']],
                data_away[['date', 'home_clean', 'away_clean', 'away_formation']],
                on=['date', 'home_clean', 'away_clean'], 
                how="inner")

# Create perspectives
home_perspective = data[["home_formation", "away_formation"]].copy()
home_perspective["outcome"] = data.apply(lambda row: 2 if row["gh"] > row["ga"] else 1 if row["gh"] == row["ga"] else 0, axis=1)
home_perspective.columns = ["our_formation", "opp_formation", "outcome"]

away_perspective = data[["away_formation", "home_formation"]].copy()
away_perspective["outcome"] = data.apply(lambda row: 2 if row["ga"] > row["gh"] else 1 if row["ga"] == row["gh"] else 0, axis=1)
away_perspective.columns = ["our_formation", "opp_formation", "outcome"]

full_data = pd.concat([home_perspective, away_perspective], ignore_index=True).dropna()

# Encode Formations
all_formations = pd.concat([full_data["our_formation"], full_data["opp_formation"]]).unique()
le = LabelEncoder()
le.fit(all_formations)

# Save the label encoder
with open("label_encoder.pkl", "wb") as f:
    pickle.dump(le, f)

print("Label encoder saved as label_encoder.pkl")

# One-hot encoding
def to_one_hot(indices, num_classes):
    indices = np.array(indices, dtype=int)
    return np.eye(num_classes)[indices]

num_classes = len(le.classes_)
X_opp = to_one_hot(full_data["opp_formation"].apply(lambda x: le.transform([x])[0]).values, num_classes)
X_our = to_one_hot(full_data["our_formation"].apply(lambda x: le.transform([x])[0]).values, num_classes)
X = np.concatenate([X_opp, X_our], axis=1)
y = full_data["outcome"].values.astype(np.float32)

# Split data
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

# Convert to tensors
X_train = torch.FloatTensor(X_train)
X_test = torch.FloatTensor(X_test)
y_train = torch.FloatTensor(y_train).view(-1, 1)
y_test = torch.FloatTensor(y_test).view(-1, 1)

# Define the Model
class FormationPredictor(nn.Module):
    def __init__(self, input_size):
        super(FormationPredictor, self).__init__()
        self.fc1 = nn.Linear(input_size, 64)
        self.fc2 = nn.Linear(64, 32)
        self.fc3 = nn.Linear(32, 1)
        self.relu = nn.ReLU()

    def forward(self, x):
        x = self.relu(self.fc1(x))
        x = self.relu(self.fc2(x))
        x = self.fc3(x)
        return x

print("Defining the model...")
model = FormationPredictor(input_size=2 * num_classes)
criterion = nn.MSELoss()
optimizer = optim.Adam(model.parameters(), lr=0.001)

# Save and load checkpoint functions
def save_checkpoint(epoch, model, optimizer, loss, filename="checkpoint.pth.tar"):
    state = {
        'epoch': epoch,
        'state_dict': model.state_dict(),
        'optimizer': optimizer.state_dict(),
        'loss': loss
    }
    torch.save(state, filename)

def load_checkpoint(filename="checkpoint.pth.tar"):
    checkpoint = torch.load(filename)
    model.load_state_dict(checkpoint['state_dict'])
    optimizer.load_state_dict(checkpoint['optimizer'])
    start_epoch = checkpoint['epoch'] + 1
    loss = checkpoint['loss']
    return start_epoch, loss

# Check if a checkpoint exists and load it
checkpoint_file = "checkpoint.pth.tar"
if os.path.isfile(checkpoint_file):
    print(f"Loading checkpoint '{checkpoint_file}'...")
    start_epoch, _ = load_checkpoint(checkpoint_file)
    print(f"Resuming training from epoch {start_epoch}...")
else:
    start_epoch = 0
    print("No checkpoint found. Starting training from scratch...")

# Train the Model
num_epochs = 5
batch_size = 64

print("Training the model...")
for epoch in tqdm(range(start_epoch, num_epochs), desc="Training epochs"):
    model.train()
    for i in tqdm(range(0, len(X_train), batch_size), desc=f"Epoch {epoch+1}/{num_epochs}", leave=False):
        batch_X = X_train[i:i+batch_size]
        batch_y = y_train[i:i+batch_size]
        optimizer.zero_grad()
        outputs = model(batch_X)
        loss = criterion(outputs, batch_y)
        loss.backward()
        optimizer.step()
    
    model.eval()
    with torch.no_grad():
        test_outputs = model(X_test)
        test_loss = criterion(test_outputs, y_test)
    print(f"Epoch {epoch+1}/{num_epochs}, Test Loss: {test_loss.item():.4f}")

    # Save the checkpoint
    save_checkpoint(epoch, model, optimizer, test_loss, filename=checkpoint_file)

# Save the model to ONNX format
print("Saving the model to ONNX format...")
model.eval()
dummy_input = torch.zeros(1, 2 * num_classes)
torch.onnx.export(model, dummy_input, "formation_predictor.onnx", 
                  input_names=["input"], output_names=["output"])
print("Model saved as formation_predictor.onnx")
