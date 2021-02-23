# twfinder
Find twitter users by them (name, handle, bio, location, followers, following count or verified)

By using Twitter API, The application start with one user and find in his (followers/following) any account match the search criteria and continue recursively with the matched users.

## How does it work
1. By using Twitter-API endpoints
    - tw-api/followers/ids.json
    - tw-api/friends/ids.json
    
    And collect the metadata for twitter users

2. Apply the search criteria on the coming users.
    - Any sub-route under configuration `SEARCH_CRITERIA` will consider as an *and condition*
    - Any option in the sub-routes use *or condition*

### Examples SEARCH_CRITERIA Configuration

- All Users in *Berlin* have in them *bio* (Developer, Software or Engineer)
```
    "SEARCH_CRITERIA": {
        "SEARCH_BIO_CONTEXT": [
            "Developer",
            "Software",
            "Engineer"
        ],
        "SEARCH_LOCATION_CONTEXT": [
            "Berlin"
        ]
    }
```

- All Users in *Silicon Valley* have in them *bio* (CEO or CTO) and have a *verified* account
```
    "SEARCH_CRITERIA": {
        "SEARCH_BIO_CONTEXT": [
            "CEO",
            "CTO"
        ],
        "SEARCH_LOCATION_CONTEXT": [
            "Silicon Valley"
        ],
        "VERIFIED": true
    }
```

- All Users have in them *NAME* (dr) and have more then 1000000 *FOLLOWERS*
```
    "SEARCH_CRITERIA": {
        "SEARCH_NAME_CONTEXT": [
            "dr"
        ],
        "FOLLOWERS_COUNT_BETWEEN": {
            "FROM": 1000000
        }
    }
```



## How To Use

### Windows Users 
1. Use the prebuilt binary in for Windows Users, Download the binary from [![HERE](https://github.com/tarekbadrshalaan/twfinder/blob/master/Windows_Users/twfinder.exe)](https://github.com/tarekbadrshalaan/twfinder/blob/master/Windows_Users/twfinder.exe)
2. Use the configuration file to configure the app (download from  [![HERE](https://github.com/tarekbadrshalaan/twfinder/blob/master/Windows_Users/config.json)](https://github.com/tarekbadrshalaan/twfinder/blob/master/Windows_Users/config.json)) Update it with your twitter credentials in `config.json`
3. Update the `SEARCH_CRITERIA` in with your search criteria.
4. run the application
5. the result will be in *result* directory


### Ubuntu Users 
1. Use the prebuilt binary in for Ubuntu Users, Download the binary from [![HERE](https://github.com/tarekbadrshalaan/twfinder/blob/master/Ubuntu_Users/twfinder)](https://github.com/tarekbadrshalaan/twfinder/blob/master/Ubuntu_Users/twfinder)
2. Use the configuration file to configure the app (download from  [![HERE](https://github.com/tarekbadrshalaan/twfinder/blob/master/Ubuntu_Users/config.json)](https://github.com/tarekbadrshalaan/twfinder/blob/master/Ubuntu_Users/config.json)) Update it with your twitter credentials in `config.json`
3. Update the `SEARCH_CRITERIA` in with your search criteria.
4. run the application
5. the result will be in *result* directory


### Build From Source
`twfinder build in golang, compile it with golang 13.0 or later`

1. Build the application by run `go build .`
2. Update the configuration file with your twitter credentials 
3. Update the `SEARCH_CRITERIA` in with your search criteria.
4. run the application
5. the result will be in *result* directory
