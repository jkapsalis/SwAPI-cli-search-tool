# SWAPI CLI Search Tool (Python, legacy)

The original Python prototype of this project. It is no longer maintained.

- Current version: [`Go/`](../Go)
- Why it was rewritten: [Why Go over Python](../README.md#why-go-over-python)

## Features

- **Character search:** prints the name, height, mass and birth year of the first matching character.
- **Homeworld details:** with `--world`, prints the homeworld's name, population, rotation period and orbital period, and compares its day and year length to Earth's.

## Run

Requires Python 3.

```bash
cd Python_old_Version
pip install -r requirements.txt    # only requests is actually used
python main.py search "luke sky"
python main.py search "luke sky" --world
```

## Sample output

![Character search output](https://github.com/user-attachments/assets/fad76912-52ec-4ec6-b697-c379a46cd911)

![Homeworld details output](https://github.com/user-attachments/assets/fbecc80e-d2d5-469a-8b53-91c287b86033)

## Code overview

| File | Contents |
|---|---|
| `main.py` | Command-line interface built with `argparse`: the `search` command and the `--world` flag |
| `swapi.py` | `search_character(name)` and `get_resource(url)`: HTTP calls to SWAPI |
| `utils.py` | `format_character`, `format_homeworld` and `calculate_time_ratio`: output formatting |
