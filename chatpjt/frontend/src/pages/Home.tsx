import { Box, CssBaseline } from "@mui/material";
import PrimayAppBar from "./templates/PrimaryAppBar";
import PrimaryDraw from "./templates/PrimaryDraw";
import SecondaryDraw from "./templates/SecondaryDraw";
import Main from "./templates/Main";

const Home = () => {

    return (
        <Box sx={{ display: "flex"}}>
            <CssBaseline />
            <PrimayAppBar />     
            <PrimaryDraw />
            <SecondaryDraw />
            <Main />
        </Box>
    );
};

export default Home;