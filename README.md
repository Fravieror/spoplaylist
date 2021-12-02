# Spoplaylist

API to Get the 100 most important songs of the date you want and add them to your Spotify playlist :smiley: :headphones:

End points are:


- `GET /hot-100/:date`  - Retrieve the top 100 songs of the given date
- `PUT /hot-100/:date` - Retrieve the top 100 songs from the given date and put them in your playlist. Before consuming this endpoint you should configure your Spotify credentials

> Note: `:date` is required with format **YYYY-MM-DD**.
