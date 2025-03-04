import pandas as pd
import pickle
from sklearn.preprocessing import LabelEncoder
import os
from tqdm import tqdm

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