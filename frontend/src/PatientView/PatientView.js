import React, { useState, useEffect, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { itemService } from "../services/api";
import "./PatientView.css";

// --- IMPORTS FOR TABS ---
import MedicalHistory from "./tabs/MedicalHistory";
import PhysicalExamination from "./tabs/PhysicalExamination";
import FamilyHistory from "./tabs/FamilyHistory";
import Surgical from "./tabs/SurgicalHistory";
import Immunization from "./tabs/Immunization";
import FemaleHistory from "./tabs/FemaleHistory";
import SocialHistory from "./tabs/SocialHistory";
import PertinentPhysicalExamination from "./tabs/PertinentPhysicalExamination";
import DiseaseSummary from "./tabs/DiseaseSummary";
import SurgerySummary from "./tabs/SurgerySummary";


function PatientView() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [patient, setPatient] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [activeTab, setActiveTab] = useState("Profiling");
  const [activeSubTab, setActiveSubTab] = useState("Medical");

  // Live clock (Philippine time)
  const [now, setNow] = useState(new Date());
  useEffect(() => {
    const t = setInterval(() => setNow(new Date()), 1000);
    return () => clearInterval(t);
  }, []);

  const fetchPatient = useCallback(async () => {
    try {
      setLoading(true);
      const data = await itemService.getItem(id);
      setPatient(data);
      setError("");
    } catch (err) {
      setError(
        "Failed to fetch patient: " + (err.response?.data?.error || err.message)
      );
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    fetchPatient();
  }, [fetchPatient]);

  const calculateAge = (birthdate) => {
    if (!birthdate) return "";
    const birth = new Date(birthdate);
    const now = new Date();
    let years = now.getFullYear() - birth.getFullYear();
    let months = now.getMonth() - birth.getMonth();
    let days = now.getDate() - birth.getDate();
    if (days < 0) {
      months--;
      days += 30;
    }
    if (months < 0) {
      years--;
      months += 12;
    }
    return `${years}Y${months}M${days}D`;
  };

  const formatDateTime = (dateTime) => {
    if (!dateTime) return "";
    return new Date(dateTime).toLocaleString("en-US", {
      month: "2-digit",
      day: "2-digit",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      hour12: true,
    });
  };

  const formatDate = (date) => {
    if (!date) return "";
    return new Date(date).toLocaleDateString("en-US", {
      month: "2-digit",
      day: "2-digit",
      year: "numeric",
    });
  };

  if (loading)
    return <div className="loading-container">Loading patient data...</div>;
  if (error) return <div className="error-container">{error}</div>;
  if (!patient) return <div className="error-container">Patient not found</div>;

  return (
    <div className="patient-view-container">
      {/* HEADER SECTION */}
      <div className="patient-header">
        <div className="header-icon">
          <div className="patient-icon">👤</div>
        </div>
        <div className="header-info">
          <div className="header-row">
            <div className="header-field">
              <label>Case No.:</label>
              <span>{patient.caseNo}</span>
            </div>
            <div className="header-field">
              <label>Hospital No.:</label>
              <span>{patient.hospitalNo}</span>
            </div>
            <div className="header-field">
              <label>Sex:</label>
              <span>{patient.sex}</span>
            </div>
          </div>
          <div className="header-row">
            <div className="header-field name-field">
              <label>Name:</label>
              <span>
                {`${patient.lastname || ""} ${patient.firstname || ""} ${
                  patient.middlename || ""
                }`
                  .trim()
                  .toUpperCase()}
              </span>
            </div>
            <div className="header-field">
              <label>Room:</label>
              <span>{patient.room}</span>
            </div>
            <div className="header-field">
              <label>Height:</label>
              <span>{patient.height} cm</span>
            </div>
          </div>
          <div className="header-row">
            <div className="header-field">
              <label>Birthdate:</label>
              <span>{formatDate(patient.birthdate)}</span>
            </div>
            <div className="header-field">
              <label>Admission:</label>
              <span>{formatDateTime(patient.admissionDate)}</span>
            </div>
            <div className="header-field">
              <label>Weight:</label>
              <span>{patient.weight} kg</span>
            </div>
          </div>
          <div className="header-row">
            <div className="header-field">
              <label>Age:</label>
              <span>{calculateAge(patient.birthdate)}</span>
            </div>
            <div className="header-field">
              <label>Discharged:</label>
              <span>{formatDateTime(patient.dischargeDate)}</span>
            </div>
            <div className="header-field">
              <label>Complaint:</label>
              <span>{patient.complaint}</span>
            </div>
          </div>
        </div>
        <div className="header-close">
          <button
            onClick={() => navigate("/dashboard")}
            className="close-button"
          >
            ✕ Close
          </button>
        </div>
      </div>

      {/* TABS CONTAINER */}
      <div className="tabs-container">
        <div className="main-tabs">
          {["Profiling", "SOAP", "Medicine"].map((tab) => (
            <button
              key={tab}
              className={`tab ${activeTab === tab ? "active" : ""}`}
              onClick={() => setActiveTab(tab)}
            >
              {tab}
            </button>
          ))}
        </div>

        {activeTab === "Profiling" && (
          <div className="sub-tabs">
            {[ 
              "Medical",
              "Surgery",
              "Surgery Summary",
              "Family",
              "Immunization",
              "Social History",
              "Female",
              "Patient History",
              "*Pertinent Physical Examinations",
              "*Physical Examination",
              "NCDQANS",
            ].map((tab) => (
              <button
                key={tab}
                className={`sub-tab ${activeSubTab === tab ? "active" : ""}`}
                onClick={() => setActiveSubTab(tab)}
              >
                {tab}
              </button>
            ))}
          </div>
        )}
      </div>

      {/* DYNAMIC CONTENT AREA */}
      <div className="content-area">
        {activeTab === "Profiling" && (
          <>
            {/* Renders MedicalHistory.js */}
            {activeSubTab === "Medical" && <MedicalHistory />}

            {/* Renders PertinentPhysicalExamination.js */}
            {activeSubTab === "*Pertinent Physical Examinations" && <PertinentPhysicalExamination />}

            {/* Renders PhysicalExamination.js */}
            {activeSubTab === "*Physical Examination" && (
              <PhysicalExamination />
            )}

            {/* Renders FamilyHistory.js */}
            {activeSubTab === "Family" && <FamilyHistory />}

            {/* Renders DiseaseSummary.js (labelled Patient History) */}
            {activeSubTab === "Patient History" && <DiseaseSummary />}

            {/* Renders SurgerySummary.js */}
            {activeSubTab === "Surgery Summary" && <SurgerySummary />}

            {/* Renders Surgical.js */}
            {activeSubTab === "Surgery" && <Surgical />}

            {/* Renders Immunization.js */}
            {activeSubTab === "Immunization" && <Immunization />}

            {/* Renders FemaleHistory.js */}
            {activeSubTab === "Female" && <FemaleHistory />}

            {/* Renders SocialHistory.js */}
            {activeSubTab === "Social History" && <SocialHistory />}
          </>
        )}

        {activeTab === "SOAP" && (
          <div className="soap-content">
            <h3>SOAP Notes</h3>
          </div>
        )}

        {activeTab === "Medicine" && (
          <div className="medicine-content">
            <h3>Medicine Records</h3>
          </div>
        )}
      </div>

      <div className="footer-info">
        <span>{new Intl.DateTimeFormat('en-US', { timeZone: 'Asia/Manila', month: '2-digit', day: '2-digit', year: 'numeric' }).format(now)}</span>
        <span>{new Intl.DateTimeFormat('en-US', { timeZone: 'Asia/Manila', hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: true }).format(now)}</span>
      </div>
    </div>
  );
}

export default PatientView;
