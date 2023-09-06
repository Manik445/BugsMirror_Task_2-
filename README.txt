This is a Backend API made with Golang 
As mentioned ,  not using any third party Package 


----------------------Steps to handle all the https request Using postman ------------------

1 . Registering a new user : 
URL: POST http://localhost:8080/Register
Headers: None
Body: JSON data representing the new user.
For example:
{
    "id": "1",
    "secretCode": "manik@123",
    "name": "manik",
    "email": "manik@123.com",
    "playlists": [
        {
            "id": "Blueeyes",
            "name": "honey singh",
            "songs": [
                {
                    "id": "blue",
                    "name": "honey",
                    "composers": "honey sing",
                    "musicURL": "https://www.youtube.com/watch?v=NbyHNASFi6U"
                }
            ]
        }
    ]
}

Method: POST
Click "Send" to register the new user. The server will generate a unique ID and secretCode for the user. 


--------------------------Login as User1-----------------------------

URL: POST http://localhost:8080/Login
Headers: None
Body: JSON data with the user's secretCode. 
For example:

{
    "id": "user1",
    "secretCode": "code1000",
    "name": "manik",
    "email": "manik@123.com",
    "playlists": [
        {
            "id": "Blueeyes",
            "name": "honey singh",
            "songs": [
                {
                    "id": "blue",
                    "name": "honey",
                    "composers": "honey sing",
                    "musicURL": "https://www.youtube.com/watch?v=NbyHNASFi6U"
                }
            ]
        }
    ]
}

Method: POST
Click "Send" to log in as User1. The server will return the user's details.

----------------------------------Creating a Playlist for User1------------------------------------:

URL: POST http://localhost:8080/CreatePlaylist
Headers: None
Body: JSON data representing the new playlist. 
For example:
{
    "id": "23",
    "name": "singh",
    "songs": [
        {
            "id": "d",
            "name": "honey",
            "composers": "honey sing",
            "musicURL": "https://www.youtube.com/watch?v=NbyHNASFi6U"
        }
    ]
}

Method: POST
Click "Send" to create a new playlist for User1. The server will generate a unique playlist ID.


----------------------------Adding a Song to User1's Playlist-----------------------------:

URL: POST http://localhost:8080/AddSongToPlaylist
Headers: None
Body: JSON data with the playlist ID and song details.
 For example:
{
    "PlaylistID": "playlist2",
    "song": {
        "id": "ds",
        "name": "honey",
        "composers": "test sing",
        "musicURL": "https://www.youtube.com/watch?v=NbyHNASFi6U"
    }
}

Method: POST
Click "Send" to add a song to User1's playlist. The server will return the updated playlist.


----------------------------View User1's Profile---------------------:

URL: GET http://localhost:8080/ViewProfile
Headers: None
Method: GET
Body : for example : 

{
    "id": "user1",
    "secretCode": "code1000",
    "name": "manik",
    "email": "manik@123.com",
    "playlists": [
        {
            "id": "Blueeyes",
            "name": "honey singh",
            "songs": [
                {
                    "id": "blue",
                    "name": "honey",
                    "composers": "honey sing",
                    "musicURL": "https://www.youtube.com/watch?v=NbyHNASFi6U"
                }
            ]
        }
    ]
}

Click "Send" to view User1's profile. The server will return User1's playlists.


------------------------Getting Song Details-----------------------

URL: GET http://localhost:8080/GetSongDetail?songID=<songID>
Headers: None
Method: GET
body: for example 
{
    "PlaylistID": "Blueeyes"
}

Replace <songID> with the actual song ID you want to retrieve. Click "Send" to get the details of the specified song. 
for example : 
songId : ds (in params)Delete a Song from User1's Playlist:


----------------------------Deleting a song from playlist-------------------------------
URL: DELETE http://localhost:8080/DeleteSongFromPlaylist
Headers: None
Body: JSON data with the playlist ID and song ID. For example:
{
    "PlaylistID": "Blueeyes"
}

Method: DELETE
Click "Send" to delete the specified song from User1's playlist.



