import React from 'react';
import Navbar from "./components/navbar.component";
import "bootstrap/dist/css/bootstrap.min.css";
import './App.css';
import { BrowserRouter as Router, Route } from "react-router-dom";
import HomePage from "./components/homepage.component";

function App() {
  return (
    <Router>
      <div className="container">
        <Navbar />
          <br/>
            <Route path="/" exact component={HomePage}/> 
      </div>
    </Router>
  );
}

export default App;
