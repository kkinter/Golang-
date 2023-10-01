import { Box, CssBaseline } from "@mui/material";
import PrimayAppBar from "./templates/PrimaryAppBar";
import PrimaryDraw from "./templates/PrimaryDraw";
import SecondaryDraw from "./templates/SecondaryDraw";
import Main from "./templates/Main";
import PopularChannels from "../components/PrimaryDraw/PopularChannels";
import ExploreCategories from "../components/SecondaryDraw/ExploreCategories";

const Home = () => {

    return (
        <Box sx={{ display: "flex"}}>
            <CssBaseline />
            <PrimayAppBar />     
            <PrimaryDraw>
                <PopularChannels open={false}/>
            </PrimaryDraw> 
            <SecondaryDraw>
                <ExploreCategories />
            </SecondaryDraw>
            <Main />
        </Box>
    );
};

export default Home;