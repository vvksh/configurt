# Configurt
Dynamic configuration management using github.

Configurt allows you to access your config stored in json format in your private github repo. It also allows you to set a `refreshInterval` (can be -1 no updates required)
If you change the configs, your programs will automatically get updated config values.

Example use case: Say  you have a scraping program which scrapes bunch of urls and you want to add more urls without restarting the program. You can add that list of url in a json config file.
```json
{
    "urls": ["abc.com", "xyz.com"]
}
```
and use `configurt` to access this config.

## How it works
On initialization, it uses the repo and filename info to fetch the content of config file, parses the json and stores it as a map. It then starts a background goroutine which refreshes the config valyes at regular intervals which is determined by `refreshInterval`.
The config map object has keys of type `string` and values of type `interface{}`. There are helper functions to get config value as strings or floats; otheriwse the user is expected to handle the type of value returned by doing necessary type casting.

## How to use it
- Create a private repo and store your configs in a json file, say `test.json` 
- Create a access token for your github account [here](https://github.com/settings/tokens))
- Get the module
    ```bash
    go get github.com/vvksh/configurt
    ```
- Import and use 
    ```go
    import github.com/vvksh/configurt
    

    accessToken := "{{personal_github_token}}"
    configurtClient := configurt.NewClient("username", accessToken, "configRepo", "configFileName", 5* time.Minute) //refresh interval set to 5 min

    configValue := configurtClient.Get("config_key") // value type is interface{}, you will need to do proper casting 
    // to cast to string;
    stringVal := configValue.(string)

    // for convenience, I've included handling of common types (string, float64)

    //string    
    configValueString := configurtClient.GetAsString("config_key")
    configValueStringArray := configurtClient.GetAsStringArray("config_key")

    // float32
    configValueFloat := configurtClient.GetAsFloat("config_key")
    configValueFloatArray := configurtClient.GetAsFloatArray("config_key")

    ```