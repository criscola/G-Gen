var imageWidth = 550;
var imagePath = "";
var canvas;
var fileNames;

$(document).ready(function () {
    $("#canvasContainer").hide();
    $(".generatorForm").slick({
        lazyLoad: 'ondemand', // ondemand progressive anticipated
        infinite: false,
        draggable: false,
        swipe: false,
        swipeToSlide: false,
        touchMove: false,
    });
    $('#forwardToGeneralSettings').click(function () {
        $(".generatorForm").slick("slickGoTo", 1);
    });
    $('#forwardToPrinterSettings').click(function () {
        $(".generatorForm").slick("slickGoTo", 3);
    });
    // Send parameters to server
    $('#forwardToLoadingScreen').click(function () {
        $('.progress .progress-bar').progressbar({ display_text: 'center' });
        $(".generatorForm").slick("slickGoTo", 2);

        // Get parameters from inputs
        let formData = {
            "jobId": getJobId(),
            "scaleFactor": $("input[name=scaleFactor]").val(),
            "modelThickness": $("input[name=modelThickness]").val(),
            "travelSpeed": $("input[name=travelSpeed]").val(),
        }

        // Send parameters to server
        $.ajax({
            url: "/generator/generate",
            type: "POST",
            data: formData,
            success: function (queueId) {
                if (queueId != "") {
                    setJobId(queueId);
                    let percentage = 0;
                    // Ask for percentage from server
                    let loadingBarUpdater = setInterval(function () {
                        $.ajax({
                            url: "/generator/queue/" + queueId,
                            type: "GET",
                            success: function (newPercentage) {
                                if (percentage != newPercentage) {
                                    percentage = newPercentage;
                                    // Get percentage from server and update loading bar
                                    $(".progress .progress-bar").attr("data-transitiongoal", percentage);
                                    $(".progress .progress-bar").progressbar();
                                    // If percentage is 100, wait a sec then slick to the preview slide
                                    if (percentage == 100) {
                                        async function toPreview() {
                                            clearInterval(loadingBarUpdater);
                                            await sleep(2000);
                                            $(".generatorForm").slick("slickGoTo", 5);
                                        }
                                        toPreview();
                                    }
                                }
                            },
                        })
                    }, 1500);
                }
            },
        })
    });
    $('#restartForm').click(function () {
        $(".generatorForm").slick("slickGoTo", 1);
    });
    $('#makeAnotherModel').click(function () {
        // Cancel all data
        $(".generatorForm").slick("slickGoTo", 0);
    });
    $('#backToImageUpload').click(function () {
        $(".generatorForm").slick("slickGoTo", 0);
    });
    $('#backToPrinterSettings').click(function () {
        $(".generatorForm").slick("slickGoTo", 3);
    });

    $('#removeImage').click(function () {
        $.ajax({
            url: '/generator/imageRemove/'+getJobId(),
            type: 'DELETE',
            success: function (result) {
                $("#forwardToGeneralSettings").prop("disabled", true);
                $("#brightnessSlider").mbsetVal(0);
                $("#contrastSlider").mbsetVal(0);
                canvas.clear();
                $("#canvasContainer").hide();
                $("#imageUpload").show();
            }
        });
    });

    // Image download button event
    $("#download").click(function () {
        console.log("jobId: " + getJobId());
        $.ajax({
            url: "/generator/outputs/"+getJobId(),
            type: "GET",
        })
    });

    // Sliders brightness/contrast
    $("#brightnessSlider").mbSlider({
        formatValue: function (val) {
            return Number(val).toFixed(2);
        },
        onSlide: function (o) {
            if (canvas != null) {
                if (canvas.getObjects()[0] != null) {
                    let val = $(o).mbgetVal();
                    let img = canvas.getObjects()[0];

                    if (img.filters[0] != null) {
                        img.filters[0]["brightness"] = parseFloat(val);
                    } else {
                        img.filters[0] = new fabric.Image.filters.Brightness({ brightness: parseFloat(val) });
                    }

                    img.applyFilters();
                    canvas.renderAll();
                }
            }
        },
        onStop: function (o) {
            // If user didn't upload an image, reset slider value to 0
            if (canvas == null || canvas.getObjects()[0] == null) {
                $("#brightnessSlider").mbsetVal(0);
            }
        }
    });

    $("#contrastSlider").mbSlider({
        formatValue: function (val) {
            return Number(val).toFixed(2);
        },
        onSlide: function (o) {
            if (canvas != null && canvas.getObjects()[0] != null) {
                let val = $(o).mbgetVal();
                let img = canvas.getObjects()[0];

                if (img.filters[1] != null) {
                    img.filters[1]["contrast"] = parseFloat(val);
                } else {
                    img.filters[1] = new fabric.Image.filters.Contrast({ contrast: parseFloat(val) });
                }

                img.applyFilters();
                canvas.renderAll();
            }
        },
        onStop: function (o) {
            // If user didn't upload an image, reset slider value to 0
            if (canvas == null || canvas.getObjects()[0] == null) {
                $("#contrastSlider").mbsetVal(0);
            }
        }

    });

    // If cookie with imageName exists, ask server if image is still there, if yes, create canvas with it
    let cookieImage = Cookies.get("ImageName");
    if (cookieImage !== undefined) {
        console.log("imageName " + cookieImage);
        createCanvasFromImage(cookieImage);
    }
    function getJobId() {
        return Cookies.get("JobId");
    }

    function setJobId(val) {
        Cookies.set("JobId", val);
    }

    Dropzone.options.imageUpload = {
        paramName: "image", // The name that will be used to transfer the file
        maxFilesize: 20, // MB
        previewsContainer: false,
        accept: function (file, done) {
            console.log("image accepted");
            // Checks if file has an image extension. If not, decline upload
            let re = /(?:\.([^.]+))?$/;
            let ext = re.exec(file.name)[1];
            ext = ext.toLowerCase();
            if (ext == "jpeg" || ext == "png" || ext == "bmp" || ext == "jpg") {
                done();
            } else {
                alert("Accepted filetypes are JPEG, JPG, PNG, BMP.");
            }
        },
        // On image upload success
        success: function (file, response) {
            console.log("upload image success");
            createCanvasFromImage(response);
            Cookies.set("ImageName", response, { expires: 3 });
        }
    };
});

function createCanvasFromImage(imageName) {
    if (imageName !== "") {
        $("#canvasContainer").show();
        let canvasContainer;

        // If canvas wasn't create before
        if (canvas == null) {
            // Create canvas element
            canvasContainer = document.createElement("canvas");
            // Append new canvas to the canvasContainer div
            $("#canvasContainer").append(canvasContainer);
            // Create new fabric.Canvas referring to the new canvas element
            canvas = new fabric.StaticCanvas(canvasContainer, {width: 0, height: 0});
            canvas.hoverCursor = "default";
        }
        // Add image to canvas
        fabric.Image.fromURL("../uploads/" + imageName + "?new=", function (oImg) {
            let scale = imageWidth / oImg.width;

            oImg.set({
                scaleX: scale,
                scaleY: scale,
                lockMovementX: true,
                lockMovementY: true,
                selectable: false,
                hasControls: false,
                hasRotatingPoint: false,

            });

            canvas.setWidth(imageWidth);
            canvas.setHeight(oImg.height * scale);
            canvas.add(oImg);
            fileNames = imageName;
        });
        $("#removeImage").prop("disabled", false);
        $('#forwardToGeneralSettings').prop("disabled", false);
        $("#imageUpload").hide();
    }
}
function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}