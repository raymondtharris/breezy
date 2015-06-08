# Breezy
Breezy is a lightweight CMS that uses markdown to write blog posts.

The webserver is written in Go and the data is stored in a MongoDB database.

Breezy focuses on a simply one column flowing blog style. This is so the content is the center of attention for the creator and viewers of the blog.

## Features
#### Markdown
The Editor in Breezy is built to handle Markdown syntax in a minimalist way. The user can just type it all out or drag and drop elements in from the side panels of the editor.

#### Backup and Export
Breezy allows you to manually backup your posts or setup up a scheduled backup system. Users can also easily export their data from the database to store or do whatever they wish.

#### Media Library
Breezy also presents a library to store images, video, and audio files to present to users on the blog or other websites.

#### Natural Language Search
On the blog side of Breezy you can enable the search feature which gives users the ability to search through the blog using natural language to find posts and media.

### BreezyNLP
BreezyNLP is the natural language processing package on Breezy. It goes through each post and parses through the markdown to determine what is important data to use when being found when searching.
