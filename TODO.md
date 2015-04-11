TODO before working version
===========================

Reading
-------

- [ ] Store unknown index fields in DB for later analysis. One field, store as text hex.
- [ ] Store unknown trackpoint fields in DB.
- [ ] Stop logging when device plugged in
- [ ] Skip too short segments while reading
- [ ] Make the too-short criteria a flag
- [ ] Clear log (with conf flag)
- [ ] SQL Indices
- [ ] Print out summaries of tracks with ID for easing upload

Uploading
---------

- [ ] Get token from Strava, store to DB
- [ ] Upload data directly from DB, based on track ids

Later
=====

- Add POIs to description, with Google links
- Fix the package naming mess
- Add upload id field to uploads table for tracking uploads
- Actually track the uploads at startup
