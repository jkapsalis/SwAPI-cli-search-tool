def format_character(char):
    return f"""
Name: {char['name']}
Height: {char['height']} cm
Mass: {char['mass']} kg
Birth Year: {char['birth_year']}
"""


def format_homeworld(world):
    return f"""
Homeworld: {world['name']}
Population: {world['population']}
Rotation Period: {world['rotation_period']} hours
Orbital Period: {world['orbital_period']} days
"""


def calculate_time_ratio(world):
    try:
        rotation = float(world["rotation_period"])
        orbital = float(world["orbital_period"])

        earth_day = 24
        earth_year = 365

        day_ratio = rotation / earth_day
        year_ratio = orbital / earth_year

        return f"""
Day vs Earth: {day_ratio:.2f}x
Year vs Earth: {year_ratio:.2f}x
"""
    except:
        return "Unknown time ratio"
