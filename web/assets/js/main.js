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
        handleFormToggle(data) {
            this.form = data || emptyForm();
            this.formOpen = !this.formOpen;
        },
        handleSubmitForm() {
            let url = "/api/quotes"; 
            let method = "POST";
            if (this.form.id) {
                url += `/${this.form.id}`;
                method = "PUT";
            }

            fetch(url, {method, body: formData(this.form)})
                .then(r => r.json())
                .then(j => {
                    if (!j.outcome) {
                        this.triggerError(j.message);
                        return
                    }
                    
                    this.triggerAlert(this.form.id ? "Quote Updated" : "Quote Added");
                    this.formOpen = false;
                    this.form = emptyForm();
                    this.getQuotes();
                })
                .catch(e => this.triggerError(`Unknown error: ${e.message}`))
        },
        async handleDelete() {
            if (!this.form.id) {
                return;
            }

            if (!confirm("Are you sure you want to delete this quote?")) {
                return;
            }

            let response = await fetch(`/api/quotes/${this.form.id}`, {method: "DELETE"})
                .then(r => r.json());

            if (!response.outcome) {
                this.triggerError("Delete Failed");
                return;
            }

            this.triggerAlert("Quote Deleted");
            this.getQuotes();
            this.handleFormToggle();
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
                    this.applyFilters();
                    this.loading = false;
                })
                .catch(() => this.triggerError("Failed to get quote list"))
        },

        filters: { idiot: "" },
        filteredList: [],
        applyFilters() {
            if (!this.filters.idiot) {
                this.filteredList = this.quoteList;
                return;
            }

            this.filteredList = [];
            for (let i in this.quoteList) {
                if (this.quoteList[i].person == this.filters.idiot) {
                    this.filteredList.push(this.quoteList[i]);
                }
            }
        },

        formatDate(date) {
            return new Date(date).toDateString();
        },
        formatText(text) {
            return text.replaceAll(/\r?\n/g, "<br />");
        }
    }))
});
