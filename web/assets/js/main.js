window.addEventListener("alpine:init", function() {
    const emptyForm = () => ({text: "", person: "", date: ""})
    const formData = obj => {
        let form = new FormData();

        for (let k in obj) {
            form.append(k, obj[k]);
        }

        return form;
    };

    Alpine.data("app", () => ({
        init() {
            this.getQuotes();
        },

        formOpen: false,
        form: emptyForm(),
        handleFormToggle() {
            this.formOpen = !this.formOpen;
        },
        handleSubmitForm() {
            fetch("/api/quotes", {method: "POST", body: formData(this.form)})
                .then(r => r.json())
                .then(j => {
                    if (!j.outcome) {
                        this.triggerError(j.message);
                        return
                    }
                    
                    this.triggerAlert("Quote Added");
                    this.formOpen = false;
                    this.form = emptyForm();
                    this.getQuotes();
                })
                .catch(e => this.triggerError(`Unknown error: ${e.message}`))
        },

        error: "",
        alert:"",
        timeout: null,
        triggerAlert(text) {
            this.resetPrompts(text, "")
            this.resetPromptTimer();
        },
        triggerError(text) {
            this.resetPrompts("", text)
            this.resetPromptTimer();
        },
        resetPromptTimer() {
            if (this.timeout !== null) {
                clearTimeout(this.timeout);
                this.timeout = null;
            }
            this.timeout = setTimeout(this.resetPrompts.bind(this), 6000);
        },
        resetPrompts(alert = "", error = "") {
            this.alert = alert;
            this.error = error;
        },

        loading: false,
        quoteList: [],
        getQuotes() {
            this.loading = true;
            fetch("/api/quotes")
                .then(r => r.json())
                .then(j => {
                    this.quoteList = j.quotes;
                    this.loading = false;
                })
                .catch(() => this.triggerError("Failed to get quote list"))
        }
    }))
});
