const progressBar = document.getElementsByClassName('progress').item(0);
const startButton = document.getElementsByClassName('start').item(0);
const taskIdTag = document.getElementsByClassName('task_id').item(0);
const estimatedTimeTag = document.getElementsByClassName('estimated_time').item(0);
const resultTag = document.getElementsByClassName('result').item(0);
const iterationsTag = document.getElementsByClassName('iterations').item(0);
const sleepTag = document.getElementsByClassName('sleep_time_ms').item(0);

let interval = null;

startButton.addEventListener('click', async () => {
    if(interval) clearInterval(interval);
    estimatedTimeTag.classList.add('hidden');
    taskIdTag.classList.add('hidden');
    resultTag.classList.add('hidden');
    progressBar.setAttribute('value', '0');

    if(+sleepTag.value <= 0) throw "wrong sleep duration";

    const data = {
        sleep_time_ms: +sleepTag.value,
        iterations: +iterationsTag.value,
    };

    const response = await postData('/long_process', data);
    if (response.status !== 200) {
        // if (response.status >= 400 && response.status < 500) {
            const error = await response.json();
            resultTag.classList.remove('hidden');
            resultTag.innerText = error.message;
        // } else {
        //     resultTag.classList.remove('hidden');
        //     resultTag.innerText = response.statusText;
        // }
        return;
    }
    const task = await response.json();
    console.log(task);
    const {id} = task;
    window.tasks = [];
    window.task = task;
    progressBar.classList.remove('hidden');
    taskIdTag.classList.remove('hidden');
    taskIdTag.innerText = `Task id: #${id}`;
    interval = setInterval(async () => {
        try {
            const response = await fetch(`/task/${id}`);
            if (response.status !== 200) {
                clearInterval(interval);
                // if (response.status >= 400 && response.status < 500) {
                    const error = await response.json();
                    resultTag.classList.remove('hidden');
                    resultTag.innerText = error.message;
                // } else {
                //     resultTag.classList.remove('hidden');
                //     resultTag.innerText = response.statusText;
                // }
                return;
            }
            const task = await response.json();
            console.log(task);
            const {done_percentage, updated_at, created_at, error, result} = task;
            if (error) {
                estimatedTimeTag.classList.add('hidden');
                resultTag.classList.remove('hidden');
                resultTag.innerHTML = error;
                clearInterval(interval);
            }
            progressBar.setAttribute('value', done_percentage);
            progressBar.innerText = `${done_percentage * 100}%`;
            if (done_percentage > 0) {
                const createdAt = new Date(created_at);
                const updatedAt = new Date(updated_at);
                const diff = updatedAt - createdAt;
                const expectedTime = new Date(createdAt - (-(diff / done_percentage)));
                const timeLeft = Math.ceil((expectedTime - updatedAt) / 1000);
                estimatedTimeTag.classList.remove('hidden');
                estimatedTimeTag.innerText = `${timeLeft} seconds left`;
            }
            if (done_percentage >= 1) {
                estimatedTimeTag.classList.add('hidden');
                if (result) {
                    resultTag.classList.remove('hidden');
                    resultTag.innerHTML = JSON.stringify(JSON.parse(atob(result)), null, 2);
                }
                clearInterval(interval);
            }
        } catch (e) {
            console.error(e);
            clearInterval(interval);
        }
    }, +sleepTag.value);
});

// Default options are marked with *
const postData = (url = '', data = {}, config = {}) => fetch(url, {
    method: "POST", // *GET, POST, PUT, DELETE, etc.
    mode: "cors", // no-cors, cors, *same-origin
    cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
    credentials: "same-origin", // include, *same-origin, omit
    headers: {
        "Content-Type": "application/json",
        // "Content-Type": "application/x-www-form-urlencoded",
    },
    redirect: "follow", // manual, *follow, error
    referrer: "no-referrer", // no-referrer, *client
    ...config,
    body: JSON.stringify(data), // body data type must match "Content-Type" header
});
