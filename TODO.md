TODO before working version
===========================

Reading
-------

- [x] Store unknown index fields in DB for later analysis. One field, store as text hex.
- [x] Store unknown trackpoint fields in DB.
- [x] Take out debug information
- [ ] Stop logging when device plugged in
- [x] Skip too short segments while reading
- [x] Make the too-short criteria a flag
- [ ] Clear log (with conf flag)
- [ ] Create SQL Indices for foreign keys
- [x] Print out summaries of tracks with ID for easing upload
- [ ] Refuse to download if track already exists

Uploading
---------

- [ ] Get token from Strava, store to DB
- [ ] Upload data directly from DB, based on track ids

Other
-----

- [ ] Godoc everything

Later
=====

- Add POIs to description, with Google links
- Fix the package naming mess
- Add upload id field to uploads table for tracking uploads
- Actually track the uploads at startup
- Tool for clearing out old debug information (can be done manually first)
- Normalize terminology (tracks vs trackpoints)
- Refactor
- Combine tracks that are close to each other

Refactoring ideas
-----------------

- CRC computation
- Concurrent track reading & database writing