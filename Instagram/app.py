import flask
from instagrapi import Client
import time
import requests
import json


def get_top_song(playlist_id, access_token):
    url = f"https://api.spotify.com/v1/playlists/{playlist_id}/tracks"
    headers = {
        "Authorization": f"Bearer {access_token}",
        "Content-Type": "application/json"
    }
    params = {
        "limit": 1,
        "fields": "items(track(name,artists(name))),total"
    }

    response = requests.get(url, headers=headers, params=params)

    if response.status_code == 200:
        data = response.json()
        top_track = data["items"][0]["track"]
        track_name = top_track["name"]
        artists = [artist["name"] for artist in top_track["artists"]]
        artist_names = ", ".join(artists)

        return (track_name, artist_names)

    else:
        print("Error:", response.status_code)


# Set your playlist ID and access token
playlist_id = "YOUR_PLAYLIST_ID"
access_token = "YOUR_ACCESS_TOKEN"

# Call the function to get the top song
get_top_song(playlist_id, access_token)


app = Flask(__name__)
USERNAME = 'sssaketh'
PASSWORD = 'Sakvith27#'

cl = Client()
cl.login(USERNAME, PASSWORD) 

authenticate()

song = ''
while True:
    time.sleep(60)

    top_song = ''
    if top_song is not song:
        song = top_song 
        cl.account_edit(biography=f'my current favorite song is {song}')





