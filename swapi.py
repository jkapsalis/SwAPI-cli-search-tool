import requests

BASE_URL = "https://swapi.dev/api"


def search_character(name):
    url = f"{BASE_URL}/people/?search={name}"
    response = requests.get(url)

    if response.status_code != 200:
        raise Exception("Failed to fetch character")

    data = response.json()
    return data["results"]


def get_resource(url):
    response = requests.get(url)

    if response.status_code != 200:
        raise Exception("Failed to fetch resource")

    return response.json()
