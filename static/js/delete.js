$('button[name="delete"]').on("click",function() {
    var name = $(this).attr("name");
    var confirmValue = confirm(id + "  kullanıcıyı silmek istiyor musunuz?");
    var id = $(this).attr("media_id");
    var url = "/api/media/" + id + "/delete";

    if (confirmValue) {
        $.ajax({
            url: url,
            method: "DELETE",
            success: function(response){
                alert(response.message)
            },
            error: function(xhr, status, error) {
                var data= xhr.responseText;
                var response= JSON.parse(data);
                alert(response);
            }
        });
    }
});


$('button[id=activate]').on("click",function() {
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