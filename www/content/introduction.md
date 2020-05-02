---
title: "Introduction"
menu: true
weight: 1
---

I did this project because I wanted to track changes to a website.
It helped me a lot, so I wanted to share my code with others.

## How it works?

Currently, web-watcher only supports static HTML pages, it analyzes the structure of the page and the tags and defines a
ratio which determines whether the page has changed or not.
Thus, text changes are not detected, only changes in the structure of the page itself are detected.

**Example**

Not detected:
```html
first.html
....
<h1>Hello world!</h1>
....

second.html
....
<h1>Goodbye world!</h1>
....
```

Detected:
```html
first.html
....
<h1>Hello world!</h1>
....

second.html
....
<h2>Goodbye world!</h2>
....
```
