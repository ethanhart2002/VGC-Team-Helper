import requests

# URL to a specific month's stats (modify as needed)
BASE_URL = "https://www.smogon.com/stats"
MONTH = "2025-03"
FILE = "gen9vgc2025regg-0"

url = f"{BASE_URL}/{MONTH}/{FILE}.txt"

response = requests.get(url)
if response.status_code == 200:
    data = response.text
    # print("Data successfully retrieved!")
else:
    print(f"Failed to retrieve data. HTTP Status: {response.status_code}")
    exit()


# Parse data (basic example for extracting Pokémon usage)
pokemon_usage = {}
lines = data.split("\n")
arr = [x for x in range(1,52)]
i = 0
while i <= 50:
    for line in lines:
        if i > 50:
            break
        elif " | {} ".format(arr[i]) in line:
            columns = line.split("|")
            rank = int(columns[1].strip())
            pokemon = columns[2].strip()
            usage = float(columns[3].strip().replace("%", ""))
            pokemon_usage[rank] = (pokemon, usage)
            i+=1


with open(f"{FILE}-{MONTH}-Usage.txt", "w") as file:
    # Write top 50 Pokémon
    for rank, (pokemon, usage) in sorted(pokemon_usage.items())[:50]:
        s = f"{pokemon}\n"
        file.write(s)
