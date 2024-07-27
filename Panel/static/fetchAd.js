(function() {
    const publisherID = new URLSearchParams(window.location.search).get('publisherID');
    const adContainer = document.getElementById('ad-container');
    
    function fetchAd() {
        fetch(`http://adserver.lontra.tech/api/ads?publisherID=${publisherID}`)
            .then(response => response.json())
            .then(data => {
                if (data && data.length > 0) {
                    const ad = data[0];
                    const adContent = '&lt;div class="ad-container"&gt;' +
                                    '&lt;img src="' + ad.imageUrl + '" alt="' + ad.title + '"&gt;' +
                                    '&lt;p&gt;' + ad.title + '&lt;/p&gt;' +
                                    '&lt;a href="' + ad.clickUrl + '" target="_blank"&gt;Click here&lt;/a&gt;' +
                                    '&lt;/div&gt;';
                    adContainer.innerHTML = adContent;
                    fetch(ad.impressionUrl);
                }
            })
            .catch(error => console.error('Error fetching ad:', error));
    }
    
    fetchAd();
})();
