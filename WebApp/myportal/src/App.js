import React from 'react';
import Navbar from "./components/navbar.component";
import "bootstrap/dist/css/bootstrap.min.css";
import './App.css';
import { BrowserRouter as Router, Route } from "react-router-dom";
import HomePage from "./components/homepage.component";
import SelectPage from "./components/selectpage.component";

function App() {
  return (
    <Router>
      <div className="container">
        <Navbar />
          <br/>
            <Route path="/" exact component={HomePage}/> 
            <Route path="/select" exact component={SelectPage}/>
      </div>
    </Router>
  );
}

export default App;
