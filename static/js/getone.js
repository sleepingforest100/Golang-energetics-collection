function createCarousel(item, container) {
    container.text('');
    var row = $('<div class="row g-1 mx-auto my-2 justify-content-around"></div>');
    var carouselItem = createCarouselItem(item);
    row.append(carouselItem);
    container.append(row);
}
function createCarouselItem(item) {
    var carousel = $('<div class="carousel col-lg-4 col-sm-6 col-10" id="tile_'+item.EnergeticsID+'"></div>');
    var card = $('<div class="card text-center d-flex justify-content-between p-2"></div>');
    var cardHeader = $('<div class="card-header"></div>').append('<h5 class="card-title text-center" ' +
        'style="color:black">' + item.Name + ' '+ item.taste+ '</h5>');
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