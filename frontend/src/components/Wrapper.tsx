import Footer from "./Footer";
import Sidebar from "./Sidebar";

const Wrapper = ({ children }: any) => {
  return (
    <div className="flex">
      <Sidebar />
      <div className="z-40 flex-1 min-h-screen text-white mx-auto">
        <main>
          {children}
        </main>
        <Footer />
      </div>
      
      <div className='logo-pattern absolute filter-[grayscale(1)] bg-[url(src/assets/logo.png)] bg-center'/>

    </div>
  );
};

export default Wrapper;
