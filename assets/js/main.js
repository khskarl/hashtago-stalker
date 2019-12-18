
let gElemTweetList = null
let gElemHashtagList = null
let gElemSearchInput = null
let gMostRecentID = 0
let gHashtagVisibility = {}

let hashtags = []

document.addEventListener("DOMContentLoaded", function () {
  // Get main elements
  gElemTweetList = document.getElementById('tweet-list')
  gElemHashtagList = document.getElementById('hashtag-list')
  gElemSearchInput = document.getElementById('search')

  // Add callbacks
  gElemSearchInput.addEventListener("keyup", function (event) {
    event.preventDefault();
    if (event.keyCode === 13) {
      addHashtag(gElemSearchInput.value)
      gElemSearchInput.value = ''
    }
  });

  postHashtags()
  getTweets()

  setInterval(getTweets, 5000)
});


function showTweetsWithHashtag(hashtagText) {
  const tweets = [...gElemTweetList.children].filter(tweet =>
    tweet.getAttribute('data-hashtags').toLowerCase().includes(hashtagText.toLowerCase())
  );

  tweets.forEach(tweet => tweet.classList.remove('hide'))
}

function hideTweetsWithHashtag(hashtagText) {
  const tweets = [...gElemTweetList.children].filter(tweet =>
    tweet.getAttribute('data-hashtags').toLowerCase().includes(hashtagText.toLowerCase())
  );

  tweets.forEach(tweet => tweet.classList.add('hide'))
}

function getTweets() {
  url = "/tweets"

  let xhr = new XMLHttpRequest()
  xhr.onreadystatechange = function () {
    if (this.readyState == XMLHttpRequest.DONE && this.status == 200) {
      let incoming_tweets = JSON.parse(this.responseText)
      console.log(incoming_tweets)
      incoming_tweets.forEach(tweet => prependTweet(tweet))

      if (incoming_tweets.length > 0)
        gMostRecentID = BigInt(incoming_tweets[incoming_tweets.length - 1]['id'])
    }
  }

  console.log("GET tweets after ID: " + gMostRecentID)
  xhr.open("GET", url + '?id=' + gMostRecentID, true)
  xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  xhr.send(null)
}

function postHashtags() {
  url = "/hashtags"

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
  if (!anyHashtagVisible(tweetData['hashtags'])) {
    tweet.classList.add('hide')
  }
  tweet.setAttribute('data-id', tweetData['id'])
  tweet.setAttribute('data-hashtags', tweetData['hashtags'])

  let tweetHeader = document.createElement('div')
  tweetHeader.classList.add('tweet-header')

  let tweetUser = document.createElement('span')
  tweetUser.classList.add('tweet-user')
  tweetUser.textContent = tweetData['user']

  let tweetTime = document.createElement('span')
  tweetTime.classList.add('tweet-time')
  const FROM_NANO_TO_SEC = 1000
  createdAt = new Date((tweetData['created_at']) * FROM_NANO_TO_SEC)
  tweetTime.textContent = createdAt

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

  gHashtagVisibility[hashtagText.toLowerCase()] = true

  hashtags.push(hashtagText)
  appendHashtag(hashtagText)

  postHashtags()
}

function removeHashtag(hashtagText) {
  let index = hashtags.indexOf(hashtagText);
  if (index > -1) {
    hashtags.splice(index, 1);
  }

  for (let i = 0; i < gElemHashtagList.children.length; i++) {
    const entry = gElemHashtagList.children[i]

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

      gHashtagVisibility[hashtagText.toLowerCase()] = true
      showTweetsWithHashtag(hashtagText.toLowerCase())
    }
    else {
      eye.textContent = '⎼'
      hashtagEntry.classList.add('greyed')

      gHashtagVisibility[hashtagText.toLowerCase()] = false
      hideTweetsWithHashtag(hashtagText.toLowerCase())
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

function anyHashtagVisible(hashtagList) {
  let anyVisible = false
  hashtagList.forEach(hashtag => {
    if (gHashtagVisibility[hashtag]) anyVisible = true
  })

  return anyVisible
}