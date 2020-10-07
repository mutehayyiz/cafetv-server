// get all media as a listed

$("#getAll").on("click", function(){
    var url = "/api/media";
    $.ajax({    //create an ajax request to Get
        type: "GET",
        url: url,
        dataType: "json",   //expect json to be returned

        success: function(response){
            var html='';
            for ( res in response){
                html+=
                    `
                        <tr media_id=`+response[res].id+`>
                                <td><button type="button" class="btn btn-default" name="display"  category="`+response[res].category+`"  media_id="`+response[res].id+`"  media_name="`+response[res].name+`" media_type="`+response[res].mediaType+`"><span>Display</span></button></td>
                                <td>`+response[res].name+`</td>
                                <td>`+response[res].description+`</td>
                                <td>`+response[res].online+`</td>
                                <td><button type="button" class="btn btn-default" name="activate"  media_id="`+response[res].id+`"><span>Activate</span></button></td>
                                <td><button type="button" class="btn btn-default" name="delete"    media_id="`+response[res].id+`"><span>Delete</span></button></td>
                        </tr>
                        `;
            }

            html+=`<script src="/static/js/delete.js"/>`;

            $("#tbody").html(html);


        }
    });
});



// display video





/*  TODO for categories
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
                    `                    <tr class="row_video" media_id=`+response[res].id+`>
                        <td>`+response[res].name+`</td>
                        <td>`+response[res].description+`</td>
                        <td>`+response[res].online+`</td>
                        <td><button type="button" class="btn btn-default" name="activate" id="activate" media_id="`+response[res].id+`"><span>Aktif</span></button></td>
                        <td><button type="button" class="btn btn-default" name="delete"   id="delete"   media_id="`+response[res].id+`"><span>Sil</span></button></td>
                    </tr>
                    `;
                }
                 //window.history.pushState('obj', 'PageTitle', url);
                $("#tbody").html(html);
            }
    });
});
*/