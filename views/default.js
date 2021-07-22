closeButton = document.getElementById('closeButton')
closeButton.addEventListener('click', async () => {
    await close(); // Call Go function
});

closeButton = document.getElementById('openFileButton')
closeButton.addEventListener('click', async () => {
    data = await openFile()
    dataDiv = document.getElementById('data')
    dataDiv.innerHTML = data
});




