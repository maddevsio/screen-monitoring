function getData() {
  var xhr = new XMLHttpRequest();
  xhr.open('GET', '/counters', true);
  xhr.timeout = 200;

  xhr.onload = function() {
    var jsonResponse = JSON.parse(this.responseText);
    username = jsonResponse.username;
    followers = jsonResponse.followed_by;
    following = jsonResponse.follows;
    posts = jsonResponse.media;
    created = jsonResponse.created;

    document.querySelector(".username").innerHTML = "Instagram: " + username;
    document.querySelector(".followers").innerHTML = "Followers: " + followers;
    document.querySelector(".following").innerHTML = "Following: " + following;
    document.querySelector(".posts").innerHTML = "Posts: " + posts;
    document.querySelector(".date").innerHTML = "Date: " + created;
  }
  xhr.onerror = function() {
    console.log('Error ' + this.status);
  }
  xhr.ontimeout = function() {
    console.log("timeout");
  };
  xhr.send(null);
};

setInterval(getData, 5000);

d3.json('counters-last-month', function(data) {
    var data_set = [
      data.media,
      data.followed_by,
      data.follows,
    ]
    for (var i = 0; i < data_set.length; i++) {
        data_set[i] = MG.convert.date(data_set[i], 'date');
    }
    MG.data_graphic({
        title: "Instagram chart",
        description: "Instagram user activity for last month.",
        data: data_set,
        width: 600,
        height: 200,
        right: 40,
        target: '#instagram-metrics',
        legend: ['Posts', 'Followers','Following'],
        legend_target: '.legend'
    });
});
