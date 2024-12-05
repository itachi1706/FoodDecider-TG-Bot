import React, {useEffect, useMemo, useState} from "react";
import {dataStats} from "@/types/dataStats";

const dataStatsList = [
  {
    icon: (
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" stroke="white" strokeWidth="3" viewBox="0 0 122.88 122.88"
           className="size-8">
        <path
          d="M61.44,0A61.44,61.44,0,1,1,0,61.44,61.44,61.44,0,0,1,61.44,0Zm-10,58.05c3.28-2.23,4.93-5.16,4.65-11.85V29c0-2.4-4.4-2.69-4.62,0l-.16,14a2,2,0,1,1-3.93,0l.17-14.44c0-2.58-4.21-2.84-4.27,0,0,4-.16,10.43-.16,14.44a1.7,1.7,0,1,1-3.34,0L40,28.61a2.34,2.34,0,0,0-3.69-1.73c-1.55,1-1.24,3-1.3,4.65L34.43,48c.09,4.79,1.35,8.68,5.09,10.33a9.73,9.73,0,0,0,2.28.59L40.51,92.07a4.22,4.22,0,0,0,4.16,4.33h.52a4.75,4.75,0,0,0,4.67-4.87L48.73,58.91a7.17,7.17,0,0,0,2.74-.86ZM71.64,90.86,71.58,61.8c-12.65-7.31-8.62-35.46,4-35.31,15.38.18,17.2,31.73,4,35.2l1,29.32c.19,6.87-8.92,7.5-8.93-.15Z"/>
      </svg>
    ),
    color: "#3FD97F",
    title: "Total Food",
    value: "...",
    key: "foodCount"
  },
  {
    icon: (
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="1.5" stroke="white"
           className="size-6">
        <path strokeLinecap="round" strokeLinejoin="round"
              d="M2.25 7.125C2.25 6.504 2.754 6 3.375 6h6c.621 0 1.125.504 1.125 1.125v3.75c0 .621-.504 1.125-1.125 1.125h-6a1.125 1.125 0 0 1-1.125-1.125v-3.75ZM14.25 8.625c0-.621.504-1.125 1.125-1.125h5.25c.621 0 1.125.504 1.125 1.125v8.25c0 .621-.504 1.125-1.125 1.125h-5.25a1.125 1.125 0 0 1-1.125-1.125v-8.25ZM3.75 16.125c0-.621.504-1.125 1.125-1.125h5.25c.621 0 1.125.504 1.125 1.125v2.25c0 .621-.504 1.125-1.125 1.125h-5.25a1.125 1.125 0 0 1-1.125-1.125v-2.25Z"/>
      </svg>
    ),
    color: "#FF9C55",
    title: "Total Groups",
    value: "...",
    key: "groupCount"
  },
  {
    icon: (
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="1.5" stroke="white"
           className="size-6">
        <path strokeLinecap="round" strokeLinejoin="round" d="M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"/>
        <path strokeLinecap="round" strokeLinejoin="round"
              d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1 1 15 0Z"/>
      </svg>
    ),
    color: "#8155FF",
    title: "Total Locations",
    value: "...",
    key: "locationCount"
  },
  {
    icon: (
      <svg
        width="26"
        height="26"
        viewBox="0 0 26 26"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <ellipse
          cx="9.75106"
          cy="6.49984"
          rx="4.33333"
          ry="4.33333"
          fill="white"
        />
        <ellipse
          cx="9.75106"
          cy="18.4178"
          rx="7.58333"
          ry="4.33333"
          fill="white"
        />
        <path
          d="M22.7496 18.4173C22.7496 20.2123 20.5445 21.6673 17.8521 21.6673C18.6453 20.8003 19.1907 19.712 19.1907 18.4189C19.1907 17.1242 18.644 16.0349 17.8493 15.1674C20.5417 15.1674 22.7496 16.6224 22.7496 18.4173Z"
          fill="white"
        />
        <path
          d="M19.4996 6.50098C19.4996 8.2959 18.0446 9.75098 16.2496 9.75098C15.8582 9.75098 15.483 9.68179 15.1355 9.55498C15.648 8.65355 15.9407 7.61084 15.9407 6.49977C15.9407 5.38952 15.6484 4.34753 15.1366 3.44656C15.4838 3.32001 15.8587 3.25098 16.2496 3.25098C18.0446 3.25098 19.4996 4.70605 19.4996 6.50098Z"
          fill="white"
        />
      </svg>
    ),
    color: "#18BFFF",
    title: "Total Users",
    value: "...",
    key: "userCount"
  },
];

const latestData = async () => {
  const data = await fetch("/api/data/general");
  return  await data.json();
}

const DataStatsBot: React.FC<dataStats> = () => {
  const [dataStatistics, setDataStatistics] = useState(dataStatsList);

  useMemo(() => {
    latestData().then((dataStats) => {
      dataStatistics.forEach((item) => {
        item.value = dataStats[item.key];
      });
      setDataStatistics([...dataStatistics]);
    })
  }, []);

  return (
    <>
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 md:gap-6 xl:grid-cols-4 2xl:gap-7.5">
        {dataStatistics.map((item, index) => (
          <div
            key={index}
            className="rounded-[10px] bg-white p-6 shadow-1 dark:bg-gray-dark"
          >
            <div
              className="flex h-14.5 w-14.5 items-center justify-center rounded-full"
              style={{backgroundColor: item.color}}
            >
              {item.icon}
            </div>

            <div className="mt-6 flex items-end justify-between">
              <div>
                <h4 className="mb-1.5 text-heading-6 font-bold text-dark dark:text-white">
                  {item.value}
                </h4>
                <span className="text-body-sm font-medium">{item.title}</span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </>
  );
};

export default DataStatsBot;
