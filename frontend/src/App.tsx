import React from 'react';
import {LandingPage} from "./containers/Landing Page";
import {BrowserRouter, Route, Switch} from "react-router-dom";
import Login from "./components/Login";


/*
Notes: Conditional rendering for the "/" route
When the user is logged in, it should automatically bring them to the main page rather than the landing one

 */
function App() {
    return (
        <div className="App">
            <BrowserRouter>
                <div>
                    <Switch>
                        <Route exact path={"/"} component={LandingPage}/>
                        <Route exact path={"/login"} component={Login}/>
                    </Switch>
                </div>
            </BrowserRouter>
        </div>
    );
}

export default App;
