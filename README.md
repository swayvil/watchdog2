![watchdog2-logo](/images/watchdog2-logo.png)
# Watchdog2
Watchdog2 is a solution to easily browse snapshots taken by surveillance cameras and stored on an email server.

![screenshot](/images/screenshot.png)

# Architecture overview
![architecture-diagram](/images/architecture-diagram.png)

# How to
Firstly you need to configure the Network Video Recorder station to send notification mails. For the exemple we will used Synology Surveillance Station.

## Synology Surveillance Station
1. Edit Notification for Camera / Motion detected:\
![notification-icon](/images/notification-icon.png)

2. Configure email notifications:\
![notification-email](/images/notification-email.png)

3. Configure motion detected notification:\
![notification-motion-detected](/images/notification-motion-detected.png)

3.1 Subject:
```
%CAMERA%-%DATE%-%TIME% Motion detected
```

3.2 Content:
```
CAMERA: %CAMERA%
DATE: %DATE%
TIME: %TIME%
```

## Prepare the credentials of the mail server
You need a login + password. On Gmail create an [App password](https://support.google.com/mail/answer/185833?hl=en)

## Install Watchdog2
1. Get Watchdog2:
```
git clone https://github.com/swayvil/watchdog2.git
```

2. Create a folder to store the snapshots images:
```
mkdir /Users/xxx/watchdog2-store
```

3. Edit **.env**:
```
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
MAIL_CRAWLER_STORE=/Users/xxx/watchdog2-store
```

4. Declare the camera names by updating the values to insert in the camera table. Edit **initdb.sh**:
```
INSERT INTO camera (camera) VALUES ('Cour');
INSERT INTO camera (camera) VALUES ('Garage');
INSERT INTO camera (camera) VALUES ('Entree');
```

5. Edit **mail-crawler/config.json**:
- Set imap server connectivity information
- Importing mail start date
- Update mail object and body parsing patterns if you set a different notification content than the example

6. Build Docker images and start the containers:
```
cd watchdog2
docker-compose -f docker-compose.yml create
docker-compose -f docker-compose.yml start
```

7. Browse the snapshots:
http://localhost:8181

# Build your own secure video surveillance system
This [article](https://swayvil.medium.com/build-a-secure-video-surveillance-system-with-your-synology-nas-e7bc755ddfeb) presents how to build your own secure video surveillance system with your Synology NAS.
