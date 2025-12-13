import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { itemService } from "./services/api";
import "./PatientView.css";

function PatientView() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [patient, setPatient] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // Tab States
  const [activeTab, setActiveTab] = useState("Profiling");
  const [activeSubTab, setActiveSubTab] = useState("*Physical Examination"); // Defaulting to the requested tab for viewing
  const [activePhyExamTab, setActivePhyExamTab] = useState("Genitourinary"); // Default inner tab

  // Medical history state (Existing)
  const [medicalHistory, setMedicalHistory] = useState({
    allergy: false,
    asthma: false,
    cancer: false,
    cerebrovascularDisease: false,
    coronaryArteryDisease: false,
    diabetesMellitus: false,
    emphysema: false,
    epilepsySeizureDisorder: false,
    hepatitis: false,
    hyperlipidemia: false,
    hypertension: false,
    pepticUlcer: false,
    pneumonia: false,
    thyroidDisease: false,
    pulmonaryTuberculosis: false,
    extrapulmonaryTuberculosis: false,
    urinaryTractInfection: false,
    mentalIllness: false,
    others: false,
  });

  // Physical Examination State
  const [phyExamData, setPhyExamData] = useState({
    generalSurvey: "Awake and Alert",
    remarks: "",
    bloodType: "A+",
    // Dynamic Sections
    skin: {
      essentiallyNormal: false,
      weakPulses: false,
      clubbing: false,
      coldClammy: false,
      cyanosisMottled: false,
      edemaSwelling: false,
      decreasedMobility: false,
      paleNailbeds: false,
      poorSkinTurgor: false,
      rashesPetechiae: false,
      others: false,
      otherText: "",
    },
    heent: {
      anictericSclerae: false,
      exudates: false,
      essentiallyNormal: false,
      abnormalPupillary: false,
      cervicalLymphadenopathy: false,
      dryMucousMembrane: false,
      ictericSclerae: false,
      paleConjunctivae: false,
      sunkenEyeballs: false,
      sunkenFontanelle: false,
      intactTympanic: false,
      pupilsBrisky: false,
      tonsillopharyngeal: false,
      hypertropicTonsils: false,
      alarFlaring: false,
      nasalDischarge: false,
      auralDischarge: false,
      palpableMass: false,
      others: false,
      otherText: "",
    },
    chest: {
      symmetricalExpansion: false,
      lumpsBreast: false,
      clearBreath: false,
      retractions: false,
      cracklesRales: false,
      wheezes: false,
      essentiallyNormal: false,
      asymmetricalExpansion: false,
      decreasedBreath: false,
      enlargeNodes: false,
      others: false,
      otherText: "",
    },
    heart: {
      adynamicPrecordium: false,
      normalRate: false,
      heavesTrills: false,
      murmurs: false,
      essentiallyNormal: false,
      displacedApex: false,
      irregularRhythm: false,
      muffledSounds: false,
      pericardialBulge: false,
      others: false,
      otherText: "",
    },
    abdomen: {
      flat: false,
      hyperactiveBowel: false,
      palpableMasses: false,
      tympaniticDull: false,
      uterineContraction: false,
      flabby: false,
      globular: false,
      muscleGuarding: false,
      tenderness: false,
      palpableMass: false,
      essentiallyNormal: false,
      abdominalRigidity: false,
      abdominalTenderness: false,
      others: false,
      otherText: "",
    },
    neurological: {
      developmentalDelay: false,
      abnormalReflexes: false,
      poorMemory: false,
      poorMuscleTone: false,
      poorCoordination: false,
      seizures: false,
      normal: false,
      motorDeficit: false,
      sensoryDeficit: false,
      essentiallyNormal: false,
      abnormalGait: false,
      abnormalPosition: false,
      abnormalSensation: false,
      others: false,
      otherText: "",
    },
    genitourinary: {
      essentiallyNormal: false,
      bloodStained: false,
      cervicalDilatation: false,
      abnormalDischarge: false,
      others: false,
      otherText: "",
    },
    digitalRectal: {
      // Placeholder based on typical structure
      essentiallyNormal: false,
      others: false,
      otherText: "",
    },
  });

  useEffect(() => {
    fetchPatient();
  }, [id]);

  const fetchPatient = async () => {
    try {
      setLoading(true);
      // Mock data handling if API fails or for demo
      const data = await itemService.getItem(id).catch(() => ({
        caseNo: "C2025-05872",
        hospitalNo: "2024-004118",
        sex: "F",
        lastname: "SONSYS LN ELEVEN",
        firstname: "SONSYS FN ELEVEN",
        room: "C",
        height: "0",
        birthdate: "1974-01-12",
        admissionDate: "2025-11-13T05:25:00",
        weight: "0",
        dischargeDate: "2025-11-13T05:25:00",
        complaint: "fever",
      }));
      setPatient(data);
      setError("");
    } catch (err) {
      setError("Failed to fetch patient");
    } finally {
      setLoading(false);
    }
  };

  const calculateAge = (birthdate) => {
    if (!birthdate) return "";
    const birth = new Date(birthdate);
    const now = new Date();
    const years = now.getFullYear() - birth.getFullYear();
    return `${years}Y`; // Simplified for display
  };

  const formatDateTime = (dateTime) => {
    if (!dateTime) return "";
    return new Date(dateTime).toLocaleString("en-US");
  };

  const formatDate = (date) => {
    if (!date) return "";
    return new Date(date).toLocaleDateString("en-US");
  };

  // Helper to handle checkbox changes in Physical Exam
  const handlePhyExamChange = (section, field, value) => {
    setPhyExamData((prev) => ({
      ...prev,
      [section]: {
        ...prev[section],
        [field]: value,
      },
    }));
  };

  // Configuration for rendering checklists to reduce code duplication
  const phyExamConfig = {
    Skin: {
      key: "skin",
      items: [
        { label: "Essentially normal", field: "essentiallyNormal" },
        { label: "Weak pulses", field: "weakPulses" },
        { label: "Clubbing", field: "clubbing" },
        { label: "Cold clammy", field: "coldClammy" },
        { label: "Cyanosis/mottled skin", field: "cyanosisMottled" },
        { label: "Edema/swelling", field: "edemaSwelling" },
        { label: "Decreased mobility", field: "decreasedMobility" },
        { label: "Pale nailbeds", field: "paleNailbeds" },
        { label: "Poor skin turgor", field: "poorSkinTurgor" },
        { label: "Rashes/Petechiae", field: "rashesPetechiae" },
      ],
    },
    HEENT: {
      key: "heent",
      items: [
        { label: "Anicteric sclerae", field: "anictericSclerae" },
        { label: "Exudates", field: "exudates" },
        { label: "Essentially Normal", field: "essentiallyNormal" },
        { label: "Abnormal pupillary reaction", field: "abnormalPupillary" },
        { label: "Cervical lymphadenopathy", field: "cervicalLymphadenopathy" },
        { label: "Dry mucous membrane", field: "dryMucousMembrane" },
        { label: "Icteric sclerae", field: "ictericSclerae" },
        { label: "Pale conjunctivae", field: "paleConjunctivae" },
        { label: "Sunken eyeballs", field: "sunkenEyeballs" },
        { label: "Sunken fontanelle", field: "sunkenFontanelle" },
        { label: "Intact tympanic membrane", field: "intactTympanic" },
        { label: "Pupils brisky reactive to light", field: "pupilsBrisky" },
        { label: "Tonsillopharyngeal congestion", field: "tonsillopharyngeal" },
        { label: "Hypertropic tonsils", field: "hypertropicTonsils" },
        { label: "Alar flaring", field: "alarFlaring" },
        { label: "Nasal discharge", field: "nasalDischarge" },
        { label: "Aural discharge", field: "auralDischarge" },
        { label: "Palpable mass", field: "palpableMass" },
      ],
    },
    Chest: {
      key: "chest",
      items: [
        { label: "Symmetrical chest expansion", field: "symmetricalExpansion" },
        { label: "Lumps over breast(s)", field: "lumpsBreast" },
        { label: "Clear breath sounds", field: "clearBreath" },
        { label: "Retractions", field: "retractions" },
        { label: "Crackles/rales", field: "cracklesRales" },
        { label: "Wheezes", field: "wheezes" },
        { label: "Essentially normal", field: "essentiallyNormal" },
        {
          label: "Asymmetrical chest expansion",
          field: "asymmetricalExpansion",
        },
        { label: "Decreased breath sounds", field: "decreasedBreath" },
        { label: "Enlarge Axillary Lymph Nodes", field: "enlargeNodes" },
      ],
    },
    Heart: {
      key: "heart",
      items: [
        { label: "Adynamic precordium", field: "adynamicPrecordium" },
        { label: "Normal rate regular rhythm", field: "normalRate" },
        { label: "Heaves/trills", field: "heavesTrills" },
        { label: "Murmurs", field: "murmurs" },
        { label: "Essentially normal", field: "essentiallyNormal" },
        { label: "Displaced apex beat", field: "displacedApex" },
        { label: "Irregular rhythm", field: "irregularRhythm" },
        { label: "Muffled heart sounds", field: "muffledSounds" },
        { label: "Pericardial bulge", field: "pericardialBulge" },
      ],
    },
    Abdomen: {
      key: "abdomen",
      items: [
        { label: "Flat", field: "flat" },
        { label: "Hyperactive bowel sounds", field: "hyperactiveBowel" },
        { label: "Palpable mass(es)", field: "palpableMasses" },
        { label: "Tympanitic/dull abdomen", field: "tympaniticDull" },
        { label: "Uterine contraction", field: "uterineContraction" },
        { label: "Flabby", field: "flabby" },
        { label: "Globular", field: "globular" },
        { label: "Muscle guarding", field: "muscleGuarding" },
        { label: "Tenderness", field: "tenderness" },
        { label: "Palpable mass", field: "palpableMass" },
        { label: "Essentially normal", field: "essentiallyNormal" },
        { label: "Abdominal rigidity", field: "abdominalRigidity" },
        { label: "Abdominal tenderness", field: "abdominalTenderness" },
      ],
    },
    Neurological: {
      key: "neurological",
      items: [
        { label: "Developmental delay", field: "developmentalDelay" },
        { label: "Abnormal reflex(es)", field: "abnormalReflexes" },
        { label: "Poor/altered memory", field: "poorMemory" },
        { label: "Poor muscle tone/strength", field: "poorMuscleTone" },
        { label: "Poor coordination", field: "poorCoordination" },
        { label: "Seizures", field: "seizures" },
        { label: "Normal", field: "normal" },
        { label: "Motor Deficit", field: "motorDeficit" },
        { label: "Sensory Deficit", field: "sensoryDeficit" },
        { label: "Essentially normal", field: "essentiallyNormal" },
        { label: "Abnormal gait", field: "abnormalGait" },
        { label: "Abnormal position sense", field: "abnormalPosition" },
        { label: "Abnormal sensation", field: "abnormalSensation" },
      ],
    },
    Genitourinary: {
      key: "genitourinary",
      items: [
        { label: "Essentially normal", field: "essentiallyNormal" },
        { label: "Blood stained in exam finger", field: "bloodStained" },
        { label: "Cervical dilatation", field: "cervicalDilatation" },
        { label: "Presence of abnormal discharge", field: "abnormalDischarge" },
      ],
    },
    "Digital Rectal": {
      key: "digitalRectal",
      items: [{ label: "Essentially normal", field: "essentiallyNormal" }],
    },
  };

  if (loading) return <div className="loading-container">Loading...</div>;
  if (error) return <div className="error-container">{error}</div>;
  if (!patient) return null;

  return (
    <div className="patient-view-container">
      {/* Header Section (Same as before) */}
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
              <span>{`${patient.lastname} ${patient.firstname}`.trim()}</span>
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

      {/* Tabs */}
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
              "Family",
              "Immunization",
              "Social History",
              "Female",
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

      <div className="content-area">
        {/* PHYSICAL EXAMINATION CONTENT */}
        {activeTab === "Profiling" &&
          activeSubTab === "*Physical Examination" && (
            <div className="physical-exam-layout">
              {/* LEFT PANE: GENERAL SURVEY */}
              <div className="pe-left-pane">
                <div className="pe-section">
                  <h4>1. General Survey</h4>
                  <div className="radio-group-vertical">
                    <label>
                      <input
                        type="radio"
                        name="genSurvey"
                        value="Awake"
                        checked={
                          phyExamData.generalSurvey === "Awake and Alert"
                        }
                        onChange={() =>
                          setPhyExamData({
                            ...phyExamData,
                            generalSurvey: "Awake and Alert",
                          })
                        }
                      />
                      <span>Awake and Alert</span>
                    </label>
                    <label>
                      <input
                        type="radio"
                        name="genSurvey"
                        value="Altered"
                        checked={
                          phyExamData.generalSurvey === "Altered Sensorium"
                        }
                        onChange={() =>
                          setPhyExamData({
                            ...phyExamData,
                            generalSurvey: "Altered Sensorium",
                          })
                        }
                      />
                      <span>Altered Sensorium</span>
                    </label>
                  </div>
                </div>

                <div className="pe-section">
                  <h4>2. Remarks</h4>
                  <textarea
                    className="remarks-box"
                    value={phyExamData.remarks}
                    onChange={(e) =>
                      setPhyExamData({
                        ...phyExamData,
                        remarks: e.target.value,
                      })
                    }
                    placeholder="lorem ipsum..."
                  />
                </div>

                <div className="pe-section bloodtype-section">
                  <div className="bloodtype-badge">Bloodtype</div>
                  <h4>1. Patient Blood Type</h4>
                  <div className="radio-group-vertical compact">
                    {["A+", "B+", "AB+", "O+", "A-", "B-", "AB-", "O-"].map(
                      (bt) => (
                        <label key={bt}>
                          <input
                            type="radio"
                            name="bloodType"
                            value={bt}
                            checked={phyExamData.bloodType === bt}
                            onChange={() =>
                              setPhyExamData({ ...phyExamData, bloodType: bt })
                            }
                          />
                          <span>{bt}</span>
                        </label>
                      )
                    )}
                  </div>
                </div>
              </div>

              {/* RIGHT PANE: SPECIFIC EXAMS */}
              <div className="pe-right-pane">
                {/* Third Level Navigation */}
                <div className="pe-nav">
                  {[
                    "Skin",
                    "HEENT",
                    "Chest",
                    "Heart",
                    "Abdomen",
                    "Neurological",
                    "Digital Rectal",
                    "Genitourinary",
                  ].map((tab) => (
                    <button
                      key={tab}
                      className={`pe-nav-link ${
                        activePhyExamTab === tab ? "active" : ""
                      }`}
                      onClick={() => setActivePhyExamTab(tab)}
                    >
                      {tab}
                    </button>
                  ))}
                </div>

                {/* Dynamic Checkbox Content */}
                <div className="pe-content-details">
                  {phyExamConfig[activePhyExamTab] && (
                    <>
                      <div className="pe-header-row">
                        <h4>
                          {/* Just numbering logic if needed, e.g. "8. Genitourinary" */}{" "}
                          {activePhyExamTab}
                        </h4>
                      </div>
                      <div className="pe-grid-layout">
                        {/* Checkboxes Area */}
                        <div className="pe-checkboxes">
                          {phyExamConfig[activePhyExamTab].items.map((item) => (
                            <label key={item.field} className="checkbox-item">
                              <input
                                type="checkbox"
                                checked={
                                  phyExamData[
                                    phyExamConfig[activePhyExamTab].key
                                  ][item.field] || false
                                }
                                onChange={(e) =>
                                  handlePhyExamChange(
                                    phyExamConfig[activePhyExamTab].key,
                                    item.field,
                                    e.target.checked
                                  )
                                }
                              />
                              <span>{item.label}</span>
                            </label>
                          ))}
                          {/* "Others" Checkbox */}
                          <label className="checkbox-item">
                            <input
                              type="checkbox"
                              checked={
                                phyExamData[phyExamConfig[activePhyExamTab].key]
                                  .others || false
                              }
                              onChange={(e) =>
                                handlePhyExamChange(
                                  phyExamConfig[activePhyExamTab].key,
                                  "others",
                                  e.target.checked
                                )
                              }
                            />
                            <span>Others</span>
                          </label>
                        </div>

                        {/* Others Text Area (Right Column in Grid) */}
                        <div className="pe-others-area">
                          <textarea
                            placeholder="lorem ipsum"
                            value={
                              phyExamData[phyExamConfig[activePhyExamTab].key]
                                .otherText
                            }
                            onChange={(e) =>
                              handlePhyExamChange(
                                phyExamConfig[activePhyExamTab].key,
                                "otherText",
                                e.target.value
                              )
                            }
                          />
                        </div>
                      </div>
                    </>
                  )}
                </div>
              </div>
            </div>
          )}

        {/* Existing Medical Content Logic */}
        {activeTab === "Profiling" && activeSubTab === "Medical" && (
          <div className="medical-content">
            {/* ... (Keep existing medical content code here) ... */}
            <div className="medical-history-section">
              <h3>Medical History Specifics</h3>
              {/* Placeholder for existing code */}
            </div>
          </div>
        )}
      </div>

      <div className="footer-info">
        <span>{new Date().toLocaleDateString("en-US")}</span>
        <span>{new Date().toLocaleTimeString("en-US")}</span>
      </div>
    </div>
  );
}

export default PatientView;
