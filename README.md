![watchdog2-logo](/images/watchdog2-logo.png)
# Watchdog2
Watchdog2 is a solution to easily browse snapshots taken by surveillance cameras and stored on an mail server.

# Architecture overview
![architecture-diagram](/images/architecture-diagram.png)

# How to
Firstly you need to configure the Network Video Recorder station to send notification mails. For the exemple we will used Synology Surveillance Station.

## Synology Surveillance Station
Edit Notification for Camera / Motion detected:
![notification-icon](/images/notification-icon.png)

Configure email notifications:
![notification-email](/images/notification-email.png)

Configure motion detected notification:
![notification-motion-detected](/images/notification-motion-detected.png)

Subject:
%CAMERA%-%DATE%-%TIME% Motion detected

Content:
CAMERA: %CAMERA%
DATE: %DATE%
TIME: %TIME%

## Install Watchdog2
Get Watchdog2:
```
git clone https://github.com/swayvil/watchdog2.git
```

Create a folder to store the snapshots images:
```
mkdir /Users/xxx/watchdog2-store
```

Edit .env:
```
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
MAIL_CRAWLER_STORE=/Users/xxx/watchdog2-store
```

Declare the camera names by updating the values to insert in the camera table:
```
INSERT INTO camera (camera) VALUES ('Cour');
INSERT INTO camera (camera) VALUES ('Garage');
INSERT INTO camera (camera) VALUES ('Entree');
```

Edit mail-crawler/config.json:
- Set imap server connectivity information
- Importing mail start date
- Update mail object and body parsing patterns if you set a different notification content than the example

Build Docker images and start the containers:
```
cd watchdog2
docker-compose -f docker-compose.yml create
docker-compose -f docker-compose.yml start
```

Browse the snapshots:
http://localhost:8181