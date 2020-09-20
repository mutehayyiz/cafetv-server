

$("[name*='btn-list']").on("click",function(){
    $("button").removeClass("active");
    $(this).addClass("active")
    var category=$(this).attr("id");
    var url = "/api/media/category/"+category;
    $.ajax({    //create an ajax request to Get
            type: "GET",
            url: url,
            dataType: "json",   //expect json to be returned
            success: function(response){
                var html='';
                for ( res in response){

                    html+=
                    `<tr class="row_video" media_id=`+response[res].id+`>
                        <td>`+response[res].name+`</td>
                        <td>`+response[res].description+`</td>
                        <td>`+response[res].online+`</td>
                        <td><button type="button" class="btn btn-default" name="activate" id="activate" media_id="`+response[res].id+`"><span>Aktif</span></button></td>
                        <td><button type="button" class="btn btn-default" name="delete"   id="delete"   media_id="`+response[res].id+`"><span>Sil</span></button></td>
                    </tr>;`
                }
                html+=`<script src="/static/js/delete.js"></script>`;
                //window.history.pushState('obj', 'PageTitle', url);
                $("#tbody").html(html);
            }
    });

});

