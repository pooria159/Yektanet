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
                                    '&lt;img src="' + ad.ImagePath + '" alt="' + ad.Title + '"&gt;' +
                                    '&lt;p&gt;' + ad.Title + '&lt;/p&gt;' +
                                    '&lt;a href="' + ad.ClickLink + '" target="_blank"&gt;Click here&lt;/a&gt;' +
                                    '&lt;/div&gt;';
                    adContainer.innerHTML = adContent;
                    fetch(ad.impressionUrl);
                }
            })
            .catch(error => console.error('Error fetching ad:', error));
    }
    
    fetchAd();
})();
