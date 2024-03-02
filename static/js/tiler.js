$(document).ready(function() {
    let items = [

    ];
    var itemsPerRow = 6;
    var url = "/energetix";
    var currentPage = 1;
    var totalPages = null;
    function fetchData(url) {
        fetch(url, {
            credentials: 'include',
        })
            .then(response => response.json())
            .then(data => {
                items = data.data;
                totalPages = data.total_count;
                loadItems(items);
                if (totalPages!==null) {
                    updatePagination(totalPages);
                }
            })
            .catch(error => console.error('Fetching error:', error));
    }
    function fetchTotal(){
        fetch('http://localhost:8080/pages', {
            credentials: 'include',
        })
            .then(response => response.json())
            .then(data => {
                totalPages = data.Pages;
            })
            .catch(error => console.error('Fetching total error:', error));
    }
    async function init(){
        await fetchTotal();
        await fetchData(constructUrl());
    }
    init();
    function loadItems(items) {
        var container = $('#tiles');
        container.empty();
        createCarousels(items, container);
    }
    function constructUrl(){
        const baseUrl = url;
        const sortSelect = document.getElementById('sortingSelect');
        const filterCheckboxes = document.querySelectorAll('input[type="checkbox"]:checked');
        const filterInputboxes = document.querySelectorAll('input[type="text"]');

        let newurl = `${baseUrl}?page=${currentPage}`;
        // Add sort option to the URL
        if (sortSelect.value !== 'default') {
            let sort = sortSelect.value.toString().split('/');
            newurl += `&sort=${sort[0]}&order=${sort[1]}`;
        }

        // Add filter options to the URL
        if (filterCheckboxes.length > 0) {
            const filters = Array.from(filterCheckboxes).map(checkbox => checkbox.value);
            for (let i = 0;i<filters.length;i++){
                newurl += `&${filters[i]}`
            }
        }
        if (filterInputboxes.length > 0){
            const inputFilters = Array.from(filterInputboxes).map(inputbox => {
                const inputValue = inputbox.value.trim();
                const inputId = inputbox.dataset.id;

                if (inputValue !== "") {
                    return `${inputId}${inputValue}`;
                }
                return null;
            });
            const filteredInputValues = inputFilters.filter(value => value !== null);
            if (filteredInputValues.length > 0) {
                for (let i = 0;i<filteredInputValues.length;i++){
                    newurl += `&${filteredInputValues[i]}`
                }
            }
        }


        console.log('Fetch url: ' +newurl)
        return newurl;
    }


    document.getElementById('sortingSelect').addEventListener('change', function (){
        currentPage = 1;
        fetchData(constructUrl());
    });
    document.getElementById('gofilter').addEventListener('click', function (){
        currentPage = 1;
        fetchData(constructUrl());
    });


    function deleteEntry(id) {
        const token = getCookie('jwtToken');
        if (!token) {
            console.error('Token not found.');
            return;
        }
        if (confirm('Are you sure you want to delete this entry?')) {
            fetch(`http://localhost:8080/energetix/${id}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`
                },
            })
                .then(response => {
                    if (response.ok) {
                        location.reload();
                    } else {
                        alert('error');
                        console.error('Failed to delete entry:', response.statusText);
                    }
                })
                .catch(error => console.error('Error deleting entry:', error));
        }
    }
    function updatePagination(totalPages) {
        const paginationElement = document.getElementById('pagination');
        const paginationList = document.getElementById('pagination-list');
        paginationList.innerHTML = ''; // Clear existing pagination links

        const maxPagesToShow = 5; // Maximum number of pages to display

        // Calculate the start and end page numbers for display
        let startPage = Math.max(1, currentPage - Math.floor(maxPagesToShow / 2));
        let endPage = Math.min(totalPages, startPage + maxPagesToShow - 1);

        // Adjust startPage and endPage if the current page is near the beginning or end
        if (endPage - startPage + 1 < maxPagesToShow) {
            startPage = Math.max(1, endPage - maxPagesToShow + 1);
        }

        // Add previous button
        if (currentPage > 1) {
            const prevPageItem = document.createElement('li');
            prevPageItem.className = 'page-item';

            const prevPageLink = document.createElement('a');
            prevPageLink.className = 'page-link';
            prevPageLink.href = '#';
            prevPageLink.textContent = 'Previous';
            prevPageLink.addEventListener('click', function(event) {
                event.preventDefault();
                currentPage--;
                fetchData(constructUrl());
            });

            prevPageItem.appendChild(prevPageLink);
            paginationList.appendChild(prevPageItem);
        }

        // Add page links
        for (let i = startPage; i <= endPage; i++) {
            const pageItem = document.createElement('li');
            pageItem.className = `page-item ${i === currentPage ? 'active' : ''}`;

            const pageLink = document.createElement('a');
            pageLink.className = 'page-link';
            pageLink.href = '#';
            pageLink.textContent = i;
            pageLink.addEventListener('click', function(event) {
                event.preventDefault();
                currentPage = i;
                fetchData(constructUrl());
            });

            pageItem.appendChild(pageLink);
            paginationList.appendChild(pageItem);
        }

        // Add next button
        if (currentPage < totalPages) {
            const nextPageItem = document.createElement('li');
            nextPageItem.className = 'page-item';

            const nextPageLink = document.createElement('a');
            nextPageLink.className = 'page-link';
            nextPageLink.href = '#';
            nextPageLink.textContent = 'Next';
            nextPageLink.addEventListener('click', function(event) {
                event.preventDefault();
                currentPage++;
                fetchData(constructUrl());
            });

            nextPageItem.appendChild(nextPageLink);
            paginationList.appendChild(nextPageItem);
        }
    }
    function createCarousels(items, container) {
        var currentIndex = 0;
        while (currentIndex < items.length) {
            var row = $('<div class="row g-1 mx-auto my-2 justify-content-around"></div>');
            for (var i = 0; i < itemsPerRow && currentIndex < items.length; i++) {
                var item = items[currentIndex];
                var carouselItem = createCarouselItem(item);
                row.append(carouselItem);
                currentIndex++;
            }
            container.append(row);
        }
    }
    function createCarouselItem(item) {
        var carousel = $('<div class="carousel col-lg-4 col-sm-6 col-10" id="tile_'+item.EnergeticsID+'"></div>');
        var card = $('<div class="card text-center d-flex justify-content-between p-2"></div>');
        var cardHeader = $('<div class="card-header"></div>').append('<h5 class="card-title text-center" ' +
            'style="color:black">' + item.Name + ' '+ item.Taste+ '</h5>');
        var cardText = $('<div class="card-text"></div>').html(
            item.Description + '<br><em class="text-warning">' +
            'CAF:' + item.Composition.Caffeine + '/TAU:' + item.Composition.Taurine + '</em>');

        var cardFooter = $('<div class="card-footer"></div>');
        var infoButton = $('<button type="button" class="btn btn-primary rounded-pill" data-bs-toggle="modal" data-bs-target="#modal_' + item.EnergeticsID + '">Info</button>');
        cardFooter.append(infoButton);
        if (auth('admin')) {
            var updateButton = $('<a href="form-go.html?id=' + item.EnergeticsID + '" class="btn btn-secondary rounded-pill">Update</a>');
            cardFooter.append(updateButton);
            var deleteButton = $('<button class="btn btn-secondary rounded-pill deletebtn" data-target=' + item.EnergeticsID + '>Delete</button>');
            cardFooter.append(deleteButton);
        }

        var toImg = $('<a class="carousel-control-prev w-100 h-75 pb-25" href="#tile_'+item.EnergeticsID+'" data-bs-slide="prev" data-bs-slide-to="0"></a>');
        card.append(cardHeader, cardText, cardFooter, toImg);

        var activeItem = $('<div class="carousel-item active"></div>');
        var img = $('<img src="' + item.PictureURL + '" class="d-block img-fluid w-100" alt="Slide">');
        var toCard =$('<a class="carousel-control-prev w-100 h-75 pb-25" href="#tile_'+item.EnergeticsID+'" data-bs-slide="prev" data-bs-slide-to="1"></a>');
        activeItem.append(img,toCard);

        var carditem = $('<div class="carousel-item"></div>');
        carditem.append(card);

        var carouselInner = $('<div class="carousel-inner"></div>').append(activeItem,carditem)
        carousel.append(carouselInner);

        createInfoModal(item);

        return carousel;
    }
    $('#tiles').on('click','.deletebtn',function(){
        var id = $(this).data('target');
        deleteEntry(id);
    })
    function createInfoModal(item){
        var modal = $('<div class="modal fade" tabindex="-1" id="modal_'+item.EnergeticsID+'" aria-hidden="true">');
        var mdialog = $('<div class="modal-dialog" role="document">');
        var mcontent = $('<div class="modal-content">');
        var mheader = $('<div class="modal-header">');
        var title = $('<h1 class="modal-title fs-5 text-center" style="color:black">'+item.Name+' '+item.Taste+'</h1>');
        var close = $('<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>');
        mheader.append(title,close);

        var mbody = $('<div class="modal-body">');
        var infopara = $('<p>'+item.Description+"<br>Manufacturer: "+item.ManufacturerName+
            ', '+item.ManufactureCountry+"<br>Nutrition facts:<br>Caffeine: " +
            item.Composition.Caffeine+"<br>Taurine: "+item.Composition.Taurine+'</p>');
        mbody.append(infopara);

        mcontent.append(mheader,mbody);
        mdialog.append(mcontent);
        modal.append(mdialog);
        $('#modalbox').append(modal);
    }
});