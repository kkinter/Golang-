import { Box, CssBaseline } from "@mui/material";
import PrimayAppBar from "./templates/PrimaryAppBar";

const Home = () => {

    return (
        <Box sx={{ display: "flex"}}>
            <CssBaseline />
            <PrimayAppBar />
        </Box>
    );
};

export default Home;