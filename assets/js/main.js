
var gElemTweetList = null
var gElemHashtagList = null
var gElemSearchInput = null

let tweets = [
  {
    id: 0,
    user: 'Gopherius',
    text: "I'm creating a new unique app called Gophitter!",
    time: "01/01/2020"
  },
  {
    id: 1,
    user: 'Gopherson',
    text: "My dad, Gopher, won't ever forgive #China. \n\n #HongKong",
    time: "01/01/2019"
  },
  {
    id: 2,
    user: 'Gopherson',
    text: "Lorem ipsum ipsum ipsum ipsum ipsum ipsum ipsum",
    time: "01/01/2019"
  },
  {
    id: 3,
    user: 'Gopheynman',
    text: "“Nobody ever figures out what life is all about, and it doesn’t matter. Explore the world. Nearly everything is really interesting if you go into it deeply enough.",
    time: "01/01/2019"
  },
]

let hashtags = []

document.addEventListener("DOMContentLoaded", function () {
  getTweets()

  // Get main elements
  gElemTweetList = document.getElementById('tweet-list')
  gElemHashtagList = document.getElementById('hashtag-list')
  gElemSearchInput = document.getElementById('search')

  // Prepopulate main elements
  // tweets.forEach(tweet => appendTweet(tweet))
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
});

function getTweets() {
  url = "http://localhost:3000/tweets"

  let xhr = new XMLHttpRequest()
  xhr.onreadystatechange = function () {
    if (this.readyState == XMLHttpRequest.DONE && this.status == 200) {
      let incoming_tweets = JSON.parse(this.responseText).reverse()
      incoming_tweets.forEach(tweet => appendTweet(tweet))
    }
  }

  xhr.open("GET", url, true)
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

function appendTweet(tweetData) {
  let tweet = document.createElement('li')
  tweet.classList.add('tweet')

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

  gElemTweetList.append(tweet)
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
