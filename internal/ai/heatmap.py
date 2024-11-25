import socket
import json
import numpy as np
import matplotlib.pyplot as plt
from matplotlib.animation import FuncAnimation
import matplotlib
matplotlib.use('TkAgg')  

# Set up the client to connect to the Go server
HOST = 'localhost'  # Update to the host's IP if running on a remote machine
PORT = 5000

client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client.connect((HOST, PORT))
print("Connected to Go server.")

# Initialize the plot
fig, ax = plt.subplots()
heatmap = ax.imshow(np.zeros((5, 5)), cmap='viridis', origin='lower', vmin=-20, vmax=20)  # Adjust size dynamically if needed
cbar = plt.colorbar(heatmap)
cbar.set_label('Potential Value')
ax.set_title('Real-Time Local Grid Potential Heatmap')
ax.set_xlabel('Grid Columns')
ax.set_ylabel('Grid Rows')

def update(frame):
    try:
        # Receive data from the Go server
        data = client.recv(4096).decode('utf-8')  # Adjust buffer size if necessary
        if not data:
            print("No data received. Closing connection.")
            client.close()
            return
        
        # Split data on newlines (multiple JSON objects might be in one recv())
        messages = data.split('\n')
        messages = messages[1:-1]
        for message in messages:
            if message.strip():  # Skip empty messages
                matrix = np.array(json.loads(message))
                print(matrix)
                heatmap.set_data(matrix)  # Update heatmap with the latest matrix

    except json.JSONDecodeError as e:
        print(f"Error during update: {e}")
    except Exception as e:
        print(f"Unexpected error: {e}")

    return heatmap

# Use Matplotlib's FuncAnimation to update the heatmap dynamically
ani = FuncAnimation(fig=fig, func=update, frames=40, interval=30)  # Update interval in milliseconds
plt.show()

# Close the client connection
client.close()

