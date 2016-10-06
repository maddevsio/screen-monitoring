function getData() {
  var xhr = new XMLHttpRequest();
  xhr.open('GET', 'http://localhost:8090/counters', true);
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
}

setInterval(getData, 5000);

d3.json('static/fake_users.json', function(data) {
    //  data = MG.convert.date(data, 'year');
    for (var i = 0; i < data.length; i++) {
        data[i] = MG.convert.date(data[i], 'date');
    }

    MG.data_graphic({
        title: "Instagram chart",
        description: "Instagram user activity for last month.",
        data: data,
        width: 600,
        height: 200,
        right: 40,
        target: '#instagram-metrics',
        legend: ['Followers','Following','Posts'],
        legend_target: '.legend'
    });
});
