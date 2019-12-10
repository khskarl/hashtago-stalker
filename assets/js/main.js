
let gElemTweetList = null
let gElemHashtagList = null
let gElemSearchInput = null

let hashtags = []

document.addEventListener("DOMContentLoaded", function () {
  // Get main elements
  gElemTweetList = document.getElementById('tweet-list')
  gElemHashtagList = document.getElementById('hashtag-list')
  gElemSearchInput = document.getElementById('search')

  // Prepopulate main elements
  let prepopulation = ["CycleToWork", "CityBike"]
  prepopulation.forEach(hashtag => addHashtag(hashtag))

  // Add callbacks
  gElemSearchInput.addEventListener("keyup", function (event) {
    event.preventDefault();
    if (event.keyCode === 13) {
      addHashtag(gElemSearchInput.value)
      gElemSearchInput.value = ''
    }
  });

  getTweets()
  setInterval(getTweets, 5000)
});

function getTweets() {
  url = "http://localhost:3000/tweets"

  let xhr = new XMLHttpRequest()
  xhr.onreadystatechange = function () {
    if (this.readyState == XMLHttpRequest.DONE && this.status == 200) {
      let incoming_tweets = JSON.parse(this.responseText).reverse()
      console.log(incoming_tweets)
      incoming_tweets.forEach(tweet => prependTweet(tweet))
    }
  }

  let mostRecentTweet = gElemTweetList.childNodes[1]
  let mostRecentID = 0
  if (mostRecentTweet != null) {
    mostRecentID = mostRecentTweet.getAttribute("data-id")
  }

  console.log("GET tweets after ID: " + mostRecentID)
  xhr.open("GET", url + '?id=' + mostRecentID, true)
  xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  xhr.send(null)
}

function postHashtags() {
  url = "http://localhost:3000/hashtags"

  let xhr = new XMLHttpRequest()
  xhr.onreadystatechange = function () {
    if (this.readyState == XMLHttpRequest.DONE && this.status == 200) {
    }
  }

  xhr.open("POST", url, true)
  xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
  data = {
    hashtags: hashtags
  }
  string = JSON.stringify(data)
  console.log(string)
  xhr.send(string)
}

function prependTweet(tweetData) {
  let tweet = document.createElement('li')
  tweet.classList.add('tweet')
  tweet.setAttribute('data-id', tweetData['id'])

  let tweetHeader = document.createElement('div')
  tweetHeader.classList.add('tweet-header')

  let tweetUser = document.createElement('span')
  tweetUser.classList.add('tweet-user')
  tweetUser.textContent = tweetData['user']

  let tweetTime = document.createElement('span')
  tweetTime.classList.add('tweet-time')
  tweetTime.textContent = tweetData['created_at']

  let tweetText = document.createElement('div')
  tweetText.classList.add('tweet-text')
  tweetText.textContent = tweetData['text']

  gElemTweetList.prepend(tweet)
  tweet.append(tweetHeader)
  tweetHeader.append(tweetUser)
  tweetHeader.append(tweetTime)
  tweet.append(tweetText)
}

function addHashtag(hashtagText) {
  for (let i = 0; i < hashtags.length; i++) {
    if (hashtags[i].toLowerCase() === hashtagText.toLowerCase())
      return
  }

  hashtags.push(hashtagText)
  appendHashtag(hashtagText)

  postHashtags()
}

function removeHashtag(hashtagText) {
  let index = hashtags.indexOf(hashtagText);
  if (index > -1) {
    hashtags.splice(index, 1);
  }

  for (let i = 1; i < gElemHashtagList.childNodes.length; i++) {
    const entry = gElemHashtagList.childNodes[i]

    if (entry.getAttribute("data-hashtag") === hashtagText) {
      entry.remove()
      break
    }
  }

  postHashtags()
}

function appendHashtag(hashtagText) {
  let hashtagEntry = document.createElement('li')
  hashtagEntry.classList.add('hashtag-entry')
  hashtagEntry.setAttribute('data-hashtag', hashtagText);

  let eye = document.createElement('button')
  eye.textContent = '👁'
  eye.addEventListener("click", function () {
    const toggleToVisible = eye.textContent != '👁'

    if (toggleToVisible) {
      eye.textContent = '👁'
      hashtagEntry.classList.remove('greyed')
    }
    else {
      eye.textContent = '⎼'
      hashtagEntry.classList.add('greyed')
    }
  });

  let hashtag = document.createElement('span')
  hashtag.classList.add('hashtag')
  hashtag.textContent = '#' + hashtagText

  let cross = document.createElement('button')
  cross.textContent = '❌'
  cross.addEventListener("click", function () {
    removeHashtag(hashtagText)
  });


  hashtagEntry.append(eye)
  hashtagEntry.append(hashtag)
  hashtagEntry.append(cross)

  gElemHashtagList.append(hashtagEntry)
}
