import tkinter as tk
from tkinter import ttk
import onnxruntime as ort
import numpy as np
import pickle
import re

# Load the ONNX model
onnx_model_path = "formation_predictor.onnx"
ort_session = ort.InferenceSession(onnx_model_path)

# Function to convert input data to one-hot encoding
def to_one_hot(indices, num_classes):
    indices = np.array(indices, dtype=int)
    return np.eye(num_classes)[indices]

# Load the label encoder
def load_label_encoder():
    with open("label_encoder.pkl", "rb") as f:
        le = pickle.load(f)
    return le

le = load_label_encoder()
num_classes = len(le.classes_)

# Function to prepare input data
def prepare_input(opponent_formation, le, num_classes):
    opponent_formation = opponent_formation.strip().strip("'\"[]")  # Ensure no leading/trailing spaces, quotes, or brackets
    opp_idx = le.transform([opponent_formation])[0] if isinstance(opponent_formation, str) else opponent_formation
    opp_one_hot = to_one_hot([opp_idx], num_classes)
    return opp_one_hot

# Function to recommend formation using ONNX model
def recommend_formation_onnx(opponent_formation, ort_session, le, num_classes):
    opp_one_hot = prepare_input(opponent_formation, le, num_classes)
    
    best_formation, best_score = None, -float("inf")
    for our_idx in range(num_classes):
        our_one_hot = to_one_hot([our_idx], num_classes)
        input_vector = np.concatenate([opp_one_hot, our_one_hot], axis=1).astype(np.float32)
        
        # Run the ONNX model
        ort_inputs = {ort_session.get_inputs()[0].name: input_vector}
        ort_outs = ort_session.run(None, ort_inputs)
        score = ort_outs[0][0, 0]
        
        if score > best_score:
            best_score = score
            best_formation = le.inverse_transform([our_idx])[0]
    
    return best_formation

# Function to handle the recommend button click
def recommend():
    opponent_formation = formation_var.get().strip().strip("'\"[]")  # Ensure no leading/trailing spaces, quotes, or brackets
    
    # Validate the format of the opponent formation
    if not re.match(r'^\d+(-\d+)+$', opponent_formation):
        result_label.config(text=f"Error: Formation '{opponent_formation}' is not in the correct format (e.g., '3-4-2-1').")
        print(f"Error: Formation '{opponent_formation}' is not in the correct format (e.g., '3-4-2-1').")
        return
    
    if opponent_formation not in le.classes_:
        result_label.config(text=f"Error: Formation '{opponent_formation}' not recognized.")
        return
    
    best_formation = recommend_formation_onnx(opponent_formation, ort_session, le, num_classes)
    result_label.config(text=f"Recommended formation: {best_formation}")

# Create the main window
root = tk.Tk()
root.title("Soccer Formation Recommender")

# Create a dropdown for selecting opponent formation
formation_var = tk.StringVar()
formation_label = ttk.Label(root, text="Select opponent formation:")
formation_label.pack(pady=5)
formation_dropdown = ttk.Combobox(root, textvariable=formation_var)
formation_dropdown['values'] = le.classes_  # Use le.classes_ directly
formation_dropdown.pack(pady=5)

# Create a button to recommend formation
recommend_button = ttk.Button(root, text="Recommend Formation", command=recommend)
recommend_button.pack(pady=10)

# Create a label to display the recommended formation
result_label = ttk.Label(root, text="")
result_label.pack(pady=5)

# Run the application
root.mainloop()
