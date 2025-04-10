import torch
import torch.nn as nn
import torch.nn.functional as F
from gamestate import GameState


class FourierFeatureEncoding:
    def __init__(self, num_frequencies=6):
        self.num_frequencies = num_frequencies
        self.freq_bands = torch.linspace(1.0, 2 ** (num_frequencies - 1), num_frequencies)

    def encode(self, x):
        # Create an encoding by applying sine and cosine at each frequency
        encoding = [torch.sin(freq * x) for freq in self.freq_bands] + \
                   [torch.cos(freq * x) for freq in self.freq_bands]
        return torch.cat(encoding, dim=-1)


class SeeGoalsDNN(nn.Module):
    def __init__(self, num_frequencies=10,
                 num_players_per_team=6,
                 num_output_features=2,
                 num_hidden_layers=3,
                 hidden_layer_size=128,
                 field_hight=9500,
                 field_width=6500,
                 ):
        super().__init__()
        self.num_frequencies = num_frequencies
        self.num_players_per_team = num_players_per_team
        self.freq_bands = torch.linspace(1.0, 2 ** (num_frequencies - 1), num_frequencies, requires_grad=False)
        if field_hight > field_width:
            self.norm_factor=torch.tensor([field_hight, field_hight], requires_grad=False).float() * 0.5
        else:
            self.norm_factor=torch.tensor([field_width, field_width], requires_grad=False).float() * 0.5
            
            
        input_size = (num_players_per_team * 2 * 2) + 2  # my team + enemy team + ball position
        enriched_size = input_size * num_frequencies * 2
        self.layers = nn.Sequential(
            nn.Linear(enriched_size, hidden_layer_size),
            nn.BatchNorm1d(hidden_layer_size),
            nn.ReLU(),
            *[nn.Sequential(nn.Linear(hidden_layer_size, hidden_layer_size), nn.BatchNorm1d(hidden_layer_size), nn.ReLU()) for _ in range(num_hidden_layers)],
            nn.Linear(hidden_layer_size, num_players_per_team * num_output_features),
            # nn.Tanh()
        )
        self.path = f"see_goals_dnn_{num_frequencies}_{num_players_per_team}_{num_output_features}_{num_hidden_layers}.pth"


    def fourier_encode(self, x):
        """Encode the input data using Fourier features"""
        encoding = [torch.sin(freq * x) for freq in self.freq_bands] + \
                   [torch.cos(freq * x) for freq in self.freq_bands]
        return torch.cat(encoding, dim=-1)

    def forward(self, my_team, enemy_team, ball):
        # Normalize the positions and only use x and y
        my_team = my_team[:, :, :2] / self.norm_factor
        enemy_team = enemy_team[:, :, :2] / self.norm_factor
        ball = ball[:, :, :2] / self.norm_factor
        
        # Concatenate the input data
        input_data = torch.cat([my_team, enemy_team, ball], dim=1)
        input_data = input_data.flatten(start_dim=1)
        
        # Fourier encode the input data
        input_data = self.fourier_encode(input_data)
        # Pass the input data through the layers
        output = self.layers(input_data)
        # Reshape the output to have the same shape as the input
        output = output.view(my_team.shape[0], self.num_players_per_team, -1)
        return output / torch.norm(output, dim=2).unsqueeze(2)


if __name__ == '__main__':

    # Instantiate model, optimizer, and loss function
    model = SeeGoalsDNN(num_players_per_team=5)

    # Example usage with GameState input
    game_state_data = {
        "RobotPositions": [
            {"Id": 1, "Team": 1, "X": 1000, "Y": 2000, "Angle": 0, "VelX": 0, "VelY": 0, "VelAngle": 0},
            {"Id": 2, "Team": 1, "X": 1500, "Y": 2500, "Angle": 0, "VelX": 0, "VelY": 0, "VelAngle": 0},
            {"Id": 3, "Team": 1, "X": 2000, "Y": 3000, "Angle": 0, "VelX": 0, "VelY": 0, "VelAngle": 0},
            {"Id": 4, "Team": 1, "X": 2500, "Y": 3500, "Angle": 0, "VelX": 0, "VelY": 0, "VelAngle": 0},
            {"Id": 5, "Team": 1, "X": 3000, "Y": 4000, "Angle": 0, "VelX": 0, "VelY": 0, "VelAngle": 0},
            {"Id": 1, "Team": 2, "X": -1000, "Y": -2000, "Angle": 0, "VelX": 0, "VelY": 0, "VelAngle": 0},
            {"Id": 2, "Team": 2, "X": -1500, "Y": -2500, "Angle": 0, "VelX": 0, "VelY": 0, "VelAngle": 0},
            {"Id": 3, "Team": 2, "X": -2000, "Y": -3000, "Angle": 0, "VelX": 0, "VelY": 0, "VelAngle": 0},
            {"Id": 4, "Team": 2, "X": -2500, "Y": -3500, "Angle": 0, "VelX": 0, "VelY": 0, "VelAngle": 0},
            {"Id": 5, "Team": 2, "X": -3000, "Y": -4000, "Angle": 0, "VelX": 0, "VelY": 0, "VelAngle": 0},
        ],
        "BallPosition": {
            "PosX": 0, "PosY": 0, "PosZ": 0, "VelX": 0, "VelY": 0, "VelZ": 0
        }
    }
    game_state = GameState(game_state_data)
    my_team = game_state.yellow_teams.to_torch().unsqueeze(0)
    enemy_team = game_state.blue_teams.to_torch().unsqueeze(0)
    ball = game_state.ball.to_torch().unsqueeze(0)

    print(my_team.shape)
    print(enemy_team.shape)
    print(ball.shape)

    output = model(my_team, enemy_team, ball)
    print(output.shape)