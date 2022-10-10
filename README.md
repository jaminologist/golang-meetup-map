# Golang-Meetup-Map

This repository is used to display upcoming local meetups based around Go on a map.

You can see the site up and running here: https://jaminologist.github.io/golang-meetup-map/

# NOW DEPRECATED
Turns out, this already exists on Meetup.com. [Go Meetup Map on Meetup](https://www.meetup.com/pro/go/) I was planning on using the Meetup API to automatically update this map with scheduled runs. However, I'd essentially be recreating the above. I'll keep the repo open as it was fun to make, but there won't be anymore updates.

Everything below this is how the Readme was before deprecation.

## How To Contribute?

If you have a meetup that you want to add to the Map, simply open a pull request and add it to the [meetups.csv](docs/meetups.csv) file located in the docs folder. 

Once merged into master, the meetup will be automatically added to the Map. 

The current layout is:

| Name        | Date           | Icon | Link | Latitude | Longitude | 
| :-------------: |:-------------:| :-----:| :-----:| :-----:| :-----:|
| Waterworld Gophers |2500-01-29| default.jpg| https://en.wikipedia.org/wiki/Waterworld | 24.9948056 | -71.0351806 |


- **Name:** Name of the Meetup/Event
- **Date:** When your meetup is scheduled.
- **Icon:** Your Meetup's Icon (Can be uploaded to the docs/icons directory in the same PR).
- **Link** Link to the event, or your Meetup's Home Page
- **Latitude/Longitude:** Location of the Meetup. 

(If your Meetup is online, you can set the location on the host country and add (online) to the name of the meetup)

