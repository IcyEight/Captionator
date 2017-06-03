## Inspiration
I was inspired to create Cool Captionator by working directly with a subject matter expert who was deaf. After discussing some issues he faced day to day, one of the more striking things he brought up was that deaf people have a difficult time understanding emotional tone from sub-titles or transcribed text from audio.
## What it does
Cool Captionator is a service whereby users can upload a mp3 or wav audio file and it will automatically be transcribed streamed back to the user as text. The transcripts will be highlighted based on the emotional tone of the message.
## How I built it
I built this using Go Language, as a file is uploaded go will copy its contents on the server and stream them to IBM Watson where it gets translated to text using their Speech-To-Text API. Upon receiving the transcripts we run them through another IBM Watson service to gather the emotional tone. We combine this data back onto the webpage where it gets styled based on the emotional tone of the message. Since the audio is being streamed in pieces, I was able to implement web sockets to stream the data back to the user as it is being evaluated by the Watson service. This provides a quick response no matter how big the audio file is. 
## Challenges I ran into
Occasionally IBM Watson will not return a valid prediction of the tone which causes some messages to highlight as red because it is the first option available. This is due to how the data is streamed from the audio file. I believe a fix for this can be achieved by breaking up the text from the stream into small parts or larger parts in order to get an accurate emotional reading. 

Another challenge was that IBM Watson has a file size limit and will terminate the service once it streams up to the limit. One way around this is to split the audio file into parts that can be streamed one after another, or concurrently and joined back later using timestamps.

EvoHaX won the Popular vote of all patrons at the hackathon (including the judges), an article can be found here about the winning project and event. https://technical.ly/philly/2017/05/01/evohax-hackathon/
