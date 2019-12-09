
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

let hashtags = ["CycleToWork", "CityBike"]

document.addEventListener("DOMContentLoaded", function () {
  getTweets()

  // Get main elements
  gElemTweetList = document.getElementById('tweet-list')
  gElemHashtagList = document.getElementById('hashtag-list')
  gElemSearchInput = document.getElementById('search')

  // Prepopulate main elements
  // tweets.forEach(tweet => appendTweet(tweet))
  hashtags.forEach(hashtag => appendHashtag('#' + hashtag))

  // Add callbacks
  gElemSearchInput.addEventListener("keyup", function (event) {
    event.preventDefault();
    if (event.keyCode === 13) {
      appendHashtag("#" + gElemSearchInput.value)
      gElemSearchInput.value = ''
    }
  });
});

function getTweets() {
  url = "http://localhost:3000/tweets"
  let xmlHttp = new XMLHttpRequest()
  xmlHttp.onreadystatechange = function () {
    if (this.readyState == 4 && this.status == 200) {
      let incoming_tweets = JSON.parse(this.responseText).reverse()
      incoming_tweets.forEach(tweet => appendTweet(tweet))
    }
  }
  xmlHttp.open("GET", url, true)
  xmlHttp.send(null)
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

function appendHashtag(hashtagText) {
  let hashtagEntry = document.createElement('li')
  hashtagEntry.classList.add('hashtag-entry')

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
  hashtag.textContent = hashtagText

  let cross = document.createElement('button')
  cross.textContent = '❌'
  cross.addEventListener("click", function () {
  });


  hashtagEntry.append(eye)
  hashtagEntry.append(hashtag)
  hashtagEntry.append(cross)

  gElemHashtagList.append(hashtagEntry)
}
