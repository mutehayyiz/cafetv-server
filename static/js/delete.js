
$('button[name="display"]').on("click",function() {

    var name = $(this).attr("media_name");

    var category = $(this).attr("category");
    var mediaType = $(this).attr("media_type");

    var html="";

    if(mediaType!= "video") {

        html = `
            <video width="400" controls name="video">
                <source id="source" src="http://localhost:4242/public/media/` + category + `/` + name + `" type="video/mp4">
                    Your browser does not support HTML video.
            </video>
            `;
    }else{
        html=
            `<img src="http://localhost:4242/public/media/` + category + `/` + name + `" alt="No media" >`;
    }

    $("#video_container").html(html);
});


$('button[name="delete"]').on("click",function() {
    var element=$(this);
    var name = element.attr("name");
    var id = element.attr("media_id");
    var confirmValue = confirm(id +":  medyayı silmek istiyor musunuz?");
    var url = "/api/media/" + id + "/delete";

    if (confirmValue) {
        $.ajax({
            url: url,
            method: "DELETE",
            success: function(response){
                alert(response.message);
                element.parent().parent().remove();
            },
            error: function(xhr, status, error) {
                var data= xhr.responseText;
                var response= JSON.parse(data);
                alert(response);
            }
        });
    }
});


$('button[name="activate"]').on("click",function() {
    var name = $(this).attr("name");
    var id = $(this).attr("media_id");
    var url = "/api/media/" + id + "/delete";
    alert(ok);
        $.ajax({
            url: url,
            method: "PUT",
            success: function(response){
                alert(response.message)
            },
            error: function(xhr, status, error) {
                var data= xhr.responseText;
                var response= JSON.parse(data);
                alert(response);
            }
        });

});


