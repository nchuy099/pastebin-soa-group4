package com.nchuy099.pastebin_soa.service;

import com.nchuy099.pastebin_soa.dto.projection.MonthlyStatsProjection;
import com.nchuy099.pastebin_soa.dto.response.MonthlyStatsResponse;

import java.time.YearMonth;

public interface PasteService {
    MonthlyStatsResponse getMonthlyStats(YearMonth yearMonth) throws Exception;
}
