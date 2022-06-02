# cyberbog

- client state stored in URL (json blob base64 compressed as able)
- "single page site"
  - POST to move
  - PUT to drop a file
  - GET to render interface 
- anything stored is decayed
  - an approach for images
  - an approach for text
  - an approach for sound file
  - generic binary bit rot for everything else
- cells have:
  - maybe a bog pool (file upload)
  - random item (generative image or text file) that can be carried and subsequently deposited
  - flavor (smell, sound, sight)
- moving moves to random cell, but does create a log in client state
- UI
  - a relevant image (decayed as needed)
  - log of actions taken
  - display of inventory
  - button to print page as PDF
  - if bog pool present:
    - button to dig
    - button to place something in pool

## bogdb, the apotheosis of nosql

upon insertion, a file is:

  - fragmented along a random fault line
  - decayed according to an algorithm (image, text) or bitrot (binary data)
  - stored with a timestamp (maybe just year)

upon reading, a fragment is:

  - fragmented along a random fault line
  - decayed according to an algorithm (image, text) or bitrot (binary data) scaled by time since insertion timestamp

## depletion

upon finding a bog pool, digging yields a random fragment. how to prevent refreshing a URL to dig over and over?

- limiting based on hash of URL state (ie the state must be updated by moving before that URL can dig again)
- limiting based on IP (weak)


## text version?

- plaintext only version of bogdb?
- MUD style interface (go $direction)
- text file oriented
- generate a ID for each discovered bog pool so it can only yield so much before moving on
