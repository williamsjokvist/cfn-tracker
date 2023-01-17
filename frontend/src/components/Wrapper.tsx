import Footer from "./Footer";
import Sidebar from "./Sidebar";

const Wrapper = ({ children }: any) => {
  return (
    <>
      <Sidebar />
      {children}
      <div className='logo-pattern absolute filter-[grayscale(1)] bg-[url(src/assets/logo.png)] bg-center'/>
    </>
  );
};

export default Wrapper;
