# Weather Service 
This is a simple Http server that listens on port 8080 and retreives the current weather for a location.

Example which will get the local weather conditions for Portland, Maine;
http://localhost:8080/weather?lat=43.68&long=-70.31

# Setup
You will need to your own API Key, please see use https://openweathermap.org/faq to get a free one! 
Then set the environment variable WEATHER_API_KEY

- go build main.go
- export WEATHER_API_KEY=yourkeygoeshere
- ./main