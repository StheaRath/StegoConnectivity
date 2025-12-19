import { useState } from 'react';
import { Generate, Extract, AnalyzeImage, AnalyzeText, SaveFile } from "../wailsjs/go/main/App";
import './style.css';

const styles = {
  container: { padding: '40px', fontFamily: '"Inter", "Segoe UI", sans-serif', background: '#f3f4f6', color: '#1f2937', minHeight: '100vh', display: 'flex', flexDirection: 'column', alignItems: 'center' },
  title: { fontSize: '18px', fontWeight: '800', color: '#111827', marginBottom: '30px', letterSpacing: '-0.5px', textAlign: 'center', maxWidth: '800px', lineHeight: '1.4' },
  card: { background: '#ffffff', padding: '35px', borderRadius: '16px', boxShadow: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)', width: '100%', maxWidth: '900px', border: '1px solid #e5e7eb', boxSizing: 'border-box' },
  nav: { display: 'flex', gap: '8px', marginBottom: '30px', background: '#e5e7eb', padding: '4px', borderRadius: '12px' },
  navBtn: { padding: '8px 24px', border: 'none', background: 'transparent', color: '#6b7280', borderRadius: '10px', fontWeight: '600', cursor: 'pointer', transition: 'all 0.2s', fontSize: '14px' },
  activeNavBtn: { background: '#ffffff', color: '#2563eb', boxShadow: '0 1px 3px rgba(0,0,0,0.1)' },
  row: { display: 'flex', gap: '15px', marginBottom: '20px', alignItems: 'flex-end', flexWrap: 'wrap' },
  col: { flex: 1, display: 'flex', flexDirection: 'column', gap: '6px', minWidth: '180px' },
  label: { fontSize: '11px', fontWeight: '700', color: '#6b7280', textTransform: 'uppercase', letterSpacing: '0.5px' },
  input: { padding: '10px 12px', borderRadius: '8px', border: '1px solid #d1d5db', background: '#f9fafb', fontSize: '14px', outline: 'none', transition: 'border 0.2s', color: '#111827', width: '100%', boxSizing: 'border-box' },
  checkboxContainer: { display: 'flex', alignItems: 'center', gap: '8px', padding: '8px 12px', border: '1px solid #d1d5db', borderRadius: '8px', background: '#f9fafb', cursor: 'pointer', height: '38px', boxSizing: 'border-box', whiteSpace: 'nowrap' },
  checkboxInput: { width: '16px', height: '16px', cursor: 'pointer', margin: 0 },
  checkboxLabel: { fontSize: '13px', fontWeight: '600', color: '#374151', cursor: 'pointer' },
  actionBtn: { width: '100%', padding: '14px', marginTop: '10px', background: '#2563eb', color: '#fff', border: 'none', borderRadius: '10px', fontSize: '15px', fontWeight: '700', cursor: 'pointer', boxShadow: '0 4px 6px rgba(37, 99, 235, 0.2)', transition: 'transform 0.1s' },
  previewContainer: { marginTop: '25px', padding: '15px', background: '#f9fafb', borderRadius: '12px', textAlign: 'center', border: '1px dashed #d1d5db' },
  image: { maxWidth: '100%', maxHeight: '350px', borderRadius: '8px', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' },
  gridWrapper: { display: 'flex', flexDirection: 'column', gap: '25px', marginTop: '25px' },
  gridBox: { display: 'flex', flexDirection: 'column', gap: '5px' },
  gridContainer: { background: '#111827', padding: '15px', borderRadius: '12px', overflow: 'auto', maxHeight: '400px', border: '1px solid #374151', boxShadow: 'inset 0 2px 4px 0 rgba(0, 0, 0, 0.2)' },
  grid: { whiteSpace: 'pre', fontFamily: '"Roboto Mono", monospace', fontSize: '11px', lineHeight: '11px', color: '#34d399', display: 'inline-block', minWidth: '100%' },
  dropzone: { padding: '50px', border: '2px dashed #9ca3af', borderRadius: '16px', textAlign: 'center', color: '#6b7280', cursor: 'pointer', background: '#f9fafb', transition: '0.2s' },
  resultBox: { marginTop: '25px', padding: '20px', background: '#eff6ff', borderLeft: '4px solid #2563eb', borderRadius: '8px' },
  subTabContainer: { display: 'flex', gap: '10px', marginBottom: '20px' },
  subTab: { flex: 1, padding: '10px', textAlign: 'center', cursor: 'pointer', borderRadius: '8px', fontWeight: '600', fontSize: '13px', border: '1px solid #d1d5db' },
  activeSubTab: { background: '#2563eb', color: '#fff', borderColor: '#2563eb' }
};

export default function App() {
  const [tab, setTab] = useState(1);
  const [gMode, setGMode] = useState("Long Text");
  const [gText, setGText] = useState("");
  const [gAlgo, setGAlgo] = useState("RSA");
  const [gKeySize, setGKeySize] = useState("2048");
  const [gConn, setGConn] = useState("4");
  const [gFeature, setGFeature] = useState("Core");
  const [gEnc, setGEnc] = useState(false);
  const [gPass, setGPass] = useState("");
  const [gRes, setGRes] = useState(null);
  const [vSource, setVSource] = useState("upload");
  const [vText, setVText] = useState("Sample Text");
  const [viewImg, setViewImg] = useState("");
  const [viewConn, setViewConn] = useState("4");
  const [viewFeature, setViewFeature] = useState("Core");
  const [analysis, setAnalysis] = useState(null);
  const [extImg, setExtImg] = useState("");
  const [extConn, setExtConn] = useState("4");
  const [extFeature, setExtFeature] = useState("Core");
  const [extDec, setExtDec] = useState(false);
  const [extPass, setExtPass] = useState("");
  const [extRes, setExtRes] = useState(null);

  const handleFile = (e, setFn) => {
    if (e.target.files && e.target.files[0]) {
      const r = new FileReader();
      r.onload = () => {
        setFn(r.result.split(',')[1]);
        if(setFn === setViewImg) setAnalysis(null);
      };
      r.readAsDataURL(e.target.files[0]);
    }
  };

  const doGen = async () => {
    try {
      let keyType = gMode === "Long Text" ? "Long Text" : gAlgo;
      const res = await Generate(gMode, gText, keyType, gKeySize, gConn, gFeature, gEnc, gPass);
      setGRes(res);
    } catch(e) { alert("Error: " + e); }
  };

  const doView = async () => {
    try {
      let res;
      if (vSource === "upload") {
        if(!viewImg) return alert("Upload Image first");
        res = await AnalyzeImage(viewImg, viewConn, viewFeature);
      } else {
        if(!vText) return alert("Enter text first");
        res = await AnalyzeText(vText, viewConn, viewFeature);
      }
      setAnalysis(res);
    } catch(e) { alert("Analysis Error: " + e); }
  };

  const doExt = async () => {
    if(!extImg) return alert("Please upload an image first.");
    const res = await Extract(extImg, extConn, extFeature, extDec, extPass);
    setExtRes(res);
  };

  return (
    <div style={styles.container}>
      <div style={styles.title}>A Steganography Tool with Image-based Symmetric Key Derivation and Ciphertext Retrieval using Connected-component Features (4/8/m Connectivity)</div>
      <div style={styles.nav}>{[1,2,3].map(i => <button key={i} style={tab===i?{...styles.navBtn, ...styles.activeNavBtn}:styles.navBtn} onClick={()=>setTab(i)}>{i===1?"1. Generate":i===2?"2. Visualize":"3. Extract"}</button>)}</div>
      <div style={styles.card}>
        {tab === 1 && <>
          <div style={styles.row}>
            <div style={styles.col}><label style={styles.label}>Payload Mode</label><select style={styles.input} value={gMode} onChange={e=>setGMode(e.target.value)}><option>Long Text</option><option>Asymmetric</option></select></div>
          </div>
          {gMode === "Long Text" ? (<div style={{...styles.col, marginBottom:'20px'}}><label style={styles.label}>Input Text</label><textarea style={{...styles.input, height:80, resize:'none'}} value={gText} onChange={e=>setGText(e.target.value)} placeholder="Message..." /></div>) : (
            <div style={styles.row}>
              <div style={styles.col}><label style={styles.label}>Algorithm</label><select style={styles.input} value={gAlgo} onChange={e=>{setGAlgo(e.target.value); const val = e.target.value; if(val === "RSA") setGKeySize("2048"); else if(val === "DH") setGKeySize("Group 14"); else setGKeySize("P-256");}}><option value="RSA">RSA</option><option value="ECC">ECC</option><option value="ECDSA">ECDSA</option><option value="DH">Diffie-Hellman</option></select></div>
              <div style={styles.col}><label style={styles.label}>Key Size / Group</label><select style={styles.input} value={gKeySize} onChange={e=>setGKeySize(e.target.value)}>{gAlgo === "RSA" && <><option value="2048">2048-bit</option><option value="3072">3072-bit</option><option value="4096">4096-bit</option></>}{gAlgo === "DH" && <><option value="Group 14">Group 14</option><option value="Group 15">Group 15</option><option value="Group 16">Group 16</option></>}{(gAlgo === "ECC" || gAlgo === "ECDSA") && <><option value="P-256">P-256</option><option value="P-384">P-384</option><option value="P-521">P-521</option></>}</select></div>
            </div>
          )}
          <div style={styles.row}>
            <div style={{...styles.col, flex:1}}><label style={styles.label}>Connectivity</label><select style={styles.input} value={gConn} onChange={e=>setGConn(e.target.value)}><option value="4">4-Conn</option><option value="8">8-Conn</option><option value="m">m-Conn</option></select></div>
            <div style={{...styles.col, flex:1}}><label style={styles.label}>Spatial Feature</label><select style={styles.input} value={gFeature} onChange={e=>setGFeature(e.target.value)}><option value="Core">Core</option><option value="Edge">Edge</option><option value="Full">Full</option><option value="Radial">Radial</option><option value="Skeleton">Skeleton</option><option value="Texture">Texture</option></select></div>
            <div style={{...styles.col, flex:0.5, minWidth:'auto'}}><label style={styles.label}>Security</label><div style={styles.checkboxContainer} onClick={() => setGEnc(!gEnc)}><input type="checkbox" checked={gEnc} onChange={() => {}} style={styles.checkboxInput} /><span style={styles.checkboxLabel}>Encrypt</span></div></div>
          </div>
          {gEnc && <div style={{...styles.col, marginTop:10}}><label style={styles.label}>Password</label><input type="password" style={styles.input} value={gPass} onChange={e=>setGPass(e.target.value)} placeholder="Password" /></div>}
          <button style={styles.actionBtn} onClick={doGen}>Generate Image</button>
          {gRes && <div style={styles.resultBox}><p style={{margin:0, fontWeight:600}}>Status: {gRes.log}</p>{gRes.privKey && <div style={{marginTop:15}}><label style={styles.label}>Key</label><textarea style={{...styles.input, width:'100%', height:60}} readOnly value={gRes.privKey} /></div>}{gRes.image && <div style={styles.previewContainer}><img src={`data:image/png;base64,${gRes.image}`} style={styles.image} /><button style={{...styles.navBtn, background:'#6b7280', color:'white', marginTop:10}} onClick={()=>SaveFile("stego.png", gRes.image)}>Download PNG</button></div>}</div>}
        </>}
        {tab === 2 && <>
          <div style={styles.subTabContainer}><div style={vSource === 'upload' ? {...styles.subTab, ...styles.activeSubTab} : styles.subTab} onClick={()=>setVSource('upload')}>UPLOAD IMAGE</div><div style={vSource === 'text' ? {...styles.subTab, ...styles.activeSubTab} : styles.subTab} onClick={()=>setVSource('text')}>INPUT TEXT</div></div>
          {vSource === 'upload' ? (<><div style={styles.dropzone} onClick={()=>document.getElementById('vfile').click()}><p style={{fontWeight:600, margin:0}}>Click to Upload</p><input id="vfile" type="file" style={{display:'none'}} onChange={e=>handleFile(e, setViewImg)} /></div>{viewImg && <div style={styles.previewContainer}><img src={`data:image/png;base64,${viewImg}`} style={styles.image} /></div>}</>) : (<div style={styles.col}><label style={styles.label}>Input Text</label><textarea style={{...styles.input, height:80}} value={vText} onChange={e=>setVText(e.target.value)} /></div>)}
          <div style={{...styles.row, marginTop:20}}>
            <div style={{...styles.col, flex:1}}><label style={styles.label}>Connectivity</label><select style={styles.input} value={viewConn} onChange={e=>setViewConn(e.target.value)}><option value="4">4-Conn</option><option value="8">8-Conn</option><option value="m">m-Conn</option></select></div>
            <div style={{...styles.col, flex:1}}><label style={styles.label}>Spatial Feature</label><select style={styles.input} value={viewFeature} onChange={e=>setViewFeature(e.target.value)}><option value="Core">Core</option><option value="Edge">Edge</option><option value="Full">Full</option><option value="Radial">Radial</option><option value="Skeleton">Skeleton</option><option value="Texture">Texture</option></select></div>
          </div>
          <button style={styles.actionBtn} onClick={doView}>Visualize Data</button>
          {analysis && <div style={styles.gridWrapper}>
            <div style={styles.gridBox}><label style={{...styles.label, color:'#e11d48'}}>1. VISUALIZED DATA PATH</label><div style={styles.gridContainer}><div style={{...styles.grid, color:'#fb7185'}}>{analysis.dataPath}</div></div></div>
            <div style={styles.gridBox}><label style={{...styles.label, color:'#2563eb'}}>2. BLOB CENTER (Conn Map)</label><div style={styles.gridContainer}><div style={{...styles.grid, color:'#60a5fa'}}>{analysis.connMap}</div></div></div>
            <div style={styles.gridBox}><label style={{...styles.label, color:'#059669'}}>3. FULL BIT GRID</label><div style={styles.gridContainer}><div style={{...styles.grid}}>{analysis.bitGrid}</div></div></div>
          </div>}
        </>}
        {tab === 3 && <>
          <div style={styles.dropzone} onClick={()=>document.getElementById('efile').click()}><p style={{fontWeight:600, margin:0}}>Click to Upload</p><input id="efile" type="file" style={{display:'none'}} onChange={e=>handleFile(e, setExtImg)} /></div>{extImg && <div style={styles.previewContainer}><img src={`data:image/png;base64,${extImg}`} style={styles.image} /></div>}
          <div style={{...styles.row, marginTop:20}}>
            <div style={{...styles.col, flex:1}}><label style={styles.label}>Connectivity</label><select style={styles.input} value={extConn} onChange={e=>setExtConn(e.target.value)}><option value="4">4-Conn</option><option value="8">8-Conn</option><option value="m">m-Conn</option></select></div>
            <div style={{...styles.col, flex:1}}><label style={styles.label}>Spatial Feature</label><select style={styles.input} value={extFeature} onChange={e=>setExtFeature(e.target.value)}><option value="Core">Core</option><option value="Edge">Edge</option><option value="Full">Full</option><option value="Radial">Radial</option><option value="Skeleton">Skeleton</option><option value="Texture">Texture</option></select></div>
            <div style={{...styles.col, flex:0.5, minWidth:'auto'}}><label style={styles.label}>Security</label><div style={styles.checkboxContainer} onClick={() => setExtDec(!extDec)}><input type="checkbox" checked={extDec} onChange={() => {}} style={styles.checkboxInput} /><span style={styles.checkboxLabel}>Decrypt</span></div></div>
          </div>
          {extDec && <div style={{...styles.col, marginTop:10}}><label style={styles.label}>Password</label><input type="password" style={styles.input} value={extPass} onChange={e=>setExtPass(e.target.value)} /></div>}
          <button style={styles.actionBtn} onClick={doExt}>Extract Payload</button>
          {extRes && <div style={styles.resultBox}><p style={{margin:0, fontWeight:600}}>Status: {extRes.log}</p>{extRes.publicKey && <div style={{marginBottom:10}}><label style={styles.label}>Public Key</label><textarea style={{...styles.input, height:60}} readOnly value={extRes.publicKey}/></div>}<div><label style={styles.label}>Extracted Data</label><textarea style={{...styles.input, height:100, fontFamily:'monospace'}} readOnly value={extRes.payload}/></div></div>}
        </>}
      </div>
    </div>
  );
}