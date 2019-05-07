function deleteHiddenInputs() {
    var inputs = document.getElementsByTagName("input");
    var elem = null;
    for (var i = 0; i < inputs.length; i++) {
        if (inputs[i].getAttribute("type") === "hidden") {
            (elem = inputs[i]).parentNode.removeChild(elem);
        }
    }
}

function patchFormActions() {
    var forms = document.getElementsByTagName("form");
    for (var i = 0; i < forms.length; i++) {
        if (forms[i].action.indexOf("http") === 0) {
            var parts = forms[i].action.split("/");
            forms[i].action = parts.length > 3 ? ("/" + parts.slice(3, parts.length).join("/")) : "/auth";
        }

    }
}

deleteHiddenInputs();
patchFormActions();
